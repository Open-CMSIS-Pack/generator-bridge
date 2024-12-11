/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/generator"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

type BridgeParamType struct {
	BoardName         string
	BoardVendor       string
	Device            string
	Output            string
	ProjectName       string
	ProjectType       string
	ForProjectPart    string
	PairedSecurePart  string
	Compiler          string
	GeneratorMap      string
	CgenName          string
	CubeContext       string
	CubeContextFolder string
}

var watcher *fsnotify.Watcher
var running bool // true if running in wait loop waiting for .ioc file

var LogFile *os.File

func procWait(proc *os.Process) {
	if proc != nil {
		if runtime.GOOS == "windows" {
			_, err := proc.Wait()
			if err != nil {
				log.Infof("Cannot wait for CubeMX to end, err %v", err)
				return
			}
		} else {
			for {
				err := proc.Signal(syscall.Signal(0))
				if err != nil {
					log.Infoln("Cannot Signal to CubeMX, is not running")
					break
				}
				time.Sleep(time.Millisecond * 200)
			}
		}
		log.Debugln("CubeMX ended")
		if watcher != nil {
			watcher.Close()
			log.Debugln("Watcher closed")
		}
		running = false // cubeMX ended, do not wait for .ioc file anymore
	}
}

func Process(cbuildGenIdxYmlPath, outPath, cubeMxPath string, runCubeMx bool, pid int) error {
	var projectFile string

	cRoot := os.Getenv("CMSIS_COMPILER_ROOT")
	if len(cRoot) == 0 {
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		exPath := filepath.Dir(ex)
		exPath = filepath.ToSlash(exPath)
		cRoot = path.Dir(exPath)
	}
	var generatorFile string
	err := filepath.Walk(cRoot, func(path string, f fs.FileInfo, err error) error {
		if f.Mode().IsRegular() && strings.Contains(path, "global.generator.yml") {
			generatorFile = path
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(generatorFile) == 0 {
		return errors.New("config file 'global.generator.yml' not found")
	}

	var gParms generator.ParamsType
	err = ReadGeneratorYmlFile(generatorFile, &gParms)
	if err != nil {
		return err
	}

	var parms cbuild.ParamsType
	err = ReadCbuildGenIdxYmlFile(cbuildGenIdxYmlPath, "CubeMX", &parms)
	if err != nil {
		return err
	}

	var bridgeParams []BridgeParamType
	err = GetBridgeInfo(&parms, &bridgeParams)
	if err != nil {
		return err
	}

	workDir := path.Dir(cbuildGenIdxYmlPath)
	if parms.Output != "" {
		if filepath.IsAbs(parms.Output) {
			workDir = parms.Output
		} else {
			workDir = path.Join(workDir, parms.Output)
		}
	} else {
		workDir = path.Join(workDir, outPath)
	}
	workDir = filepath.Clean(workDir)
	workDir = filepath.ToSlash(workDir)

	err = os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		return err
	}

	cubeIocPath := workDir
	if pid >= 0 {
		lastPath := filepath.Base(cubeIocPath)
		if lastPath != "STM32CubeMX" {
			cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX")
		}
		iocprojectPath := path.Join(cubeIocPath, "STM32CubeMX.ioc")
		mxprojectPath := path.Join(cubeIocPath, ".mxproject")
		log.Debugf("pid of CubeMX in daemon: %d", pid)
		running = true
		first := true
		iocProjectWait := false
		for {
			proc, err := os.FindProcess(pid) // this only works for windows as it is now
			if err == nil {                  // cubeMX already runs
				if runtime.GOOS != "windows" {
					err = proc.Signal(syscall.Signal(0))
					if err != nil {
						break // out of loop if CubeMX does not run anymore
					}
				}
				if first { // only start wait thread once
					go procWait(proc)
					first = false
				}
				if !running {
					break // out of loop if CubeMX does not run anymore
				}
				stIOC, err := os.Stat(iocprojectPath)
				if err != nil { // .ioc file not (yet) there
					iocProjectWait = true
					time.Sleep(time.Second)
					continue // stay in loop waiting for CubeMX end or change of .mxproject
				}
				if iocProjectWait { // .ioc file was created new, there will not be multiple changes
					log.Debugln("new project file:", iocprojectPath)
					i := 1
					for ; i < 100; i++ { // wait for .mxproject coming
						_, err := os.Stat(mxprojectPath)
						if err == nil {
							break // .mxproject appeared
						}
						time.Sleep(time.Second)
					}
					if i < 100 {
						mxproject, err := IniReader(mxprojectPath, bridgeParams)
						if err != nil {
							continue // stay in loop waiting for CubeMX end or change of .mxproject
						}
						err = ReadContexts(iocprojectPath, bridgeParams)
						if err != nil {
							continue // stay in loop waiting for CubeMX end or change of .mxproject
						}
						err = WriteCgenYml(workDir, mxproject, bridgeParams)
						if err != nil {
							continue // stay in loop waiting for CubeMX end or change of .mxproject
						}
						iocProjectWait = false // reset wait for iocProject flag
					}
				} else {
					st, err := os.Stat(mxprojectPath)
					if err != nil {
						continue // stay in loop waiting for CubeMX end or change of .mxproject
					}
					fIoc, err := os.Open(mxprojectPath)
					if err != nil {
						continue // stay in loop waiting for CubeMX end or change of .mxproject
					}
					mxprojectBuf := make([]byte, st.Size())
					_, err = fIoc.Read(mxprojectBuf)
					if err != nil {
						log.Fatal(err)
					}
					fIoc.Close()
					for { // wait for .mxproject change
						if !running {
							break // out of loop if CubeMX does not run anymore
						}
						time.Sleep(time.Second)
						stIOC1, err := os.Stat(iocprojectPath)
						if err != nil { // .ioc file not (yet) there
							break // continue loop waiting for CubeMX end or change of .mxproject
						}
						st1, err := os.Stat(mxprojectPath)
						if err != nil {
							break // continue loop waiting for CubeMX end or change of .mxproject
						}
						if stIOC.ModTime() != stIOC1.ModTime() && st.ModTime() != st1.ModTime() { // time changed
							// it seems to me that the compare is superfluous because there are only in rare cases changes but the time always changes
							if st.Size() == st1.Size() { // no change in length, compare content
								fIoc, err := os.Open(mxprojectPath)
								if err != nil {
									break // continue loop waiting for CubeMX end or change of .mxproject
								}
								mxprojectBuf1 := make([]byte, st1.Size())
								_, err = fIoc.Read(mxprojectBuf1)
								if err != nil {
									log.Fatal(err)
								}
								fIoc.Close()
								// if bytes.Equal(mxprojectBuf, mxprojectBuf1) {
								//	continue // wait for .mxproject change
								// }
							}
							mxproject, err := IniReader(mxprojectPath, bridgeParams)
							if err != nil {
								break // continue loop waiting for CubeMX end or change of .mxproject
							}
							err = ReadContexts(iocprojectPath, bridgeParams)
							if err != nil {
								break // continue loop waiting for CubeMX end or change of .mxproject
							}
							log.Debugln("Writing Cgen.yml file")
							err = WriteCgenYml(workDir, mxproject, bridgeParams)
							if err != nil {
								break // continue loop waiting for CubeMX end or change of .mxproject
							}
							break // leave inner loop reload all
						}
					}
				}
			} else {
				break // CubeMX does not run anymore
			}
		}
		// should only come here if CubeMX does not run anymore
	}

	if runCubeMx {
		lastPath := filepath.Base(cubeIocPath)
		if lastPath != "STM32CubeMX" {
			cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX")
		}
		cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX.ioc")

		var err error
		var pid int
		if utils.FileExists(cubeIocPath) {
			pid, err = Launch(cubeIocPath, "")
			if err != nil {
				return errors.New("generator '" + gParms.ID + "' missing. Install from '" + gParms.DownloadURL + "'")
			}
		} else {
			projectFile, err = WriteProjectFile(workDir, bridgeParams[0])
			if err != nil {
				return err
			}
			log.Debugf("Generated file: %v", projectFile)

			pid, err = Launch("", projectFile)
			if err != nil {
				return errors.New("generator '" + gParms.ID + "' missing. Install from '" + gParms.DownloadURL + "'")
			}
		}
		// here cubeMX runs
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		ownPath, err := filepath.EvalSymlinks((exe))
		if err != nil {
			return err
		}
		cmd := exec.Command(ownPath) //nolint
		cmd.Args = os.Args
		cmd.Args = append(cmd.Args, "-p", fmt.Sprint(pid)) // pid of cubeMX
		log.Debugf("cmd.Start as %v", cmd)
		if err := cmd.Start(); err != nil { // start myself as a daemon
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func Launch(iocFile, projectFile string) (int, error) {
	const cubeEnvVar = "STM32CubeMX_PATH"
	cubeEnv := os.Getenv(cubeEnvVar)
	if cubeEnv == "" {
		return -1, errors.New("environment variable for CubeMX not set: " + cubeEnvVar)
	}

	if iocFile != "" {
		log.Infoln("Launching STM32CubeMX with ", iocFile)
	} else if projectFile != "" {
		log.Infoln("Launching STM32CubeMX with -s ", projectFile)
	} else {
		log.Infoln("Launching STM32CubeMX...")
	}

	var pathJava string
	var arg0 string
	var arg1 string
	var pathCubeMx string
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		pathJava = path.Join(cubeEnv, "jre", "bin", "java.exe")
		pathCubeMx = path.Join(cubeEnv, "STM32CubeMX.exe")
	case "darwin":
		pathJava = path.Join(cubeEnv, "jre", "Contents", "Home", "bin", "java")
		arg0 = path.Join(cubeEnv, "stm32cubemx.icns")
		arg0 = "-Xdock:icon=" + arg0
		arg1 = "-Xdock:name=STM32CubeMX"
		pathCubeMx = path.Join(cubeEnv, "STM32CubeMX")
	default:
		pathJava = path.Join(cubeEnv, "jre", "bin", "java")
		pathCubeMx = path.Join(cubeEnv, "STM32CubeMX")
	}

	if runtime.GOOS == "darwin" {
		if iocFile != "" {
			cmd = exec.Command(pathJava, arg0, arg1, "-jar", pathCubeMx, iocFile)
		} else if projectFile != "" {
			cmd = exec.Command(pathJava, arg0, arg1, "-jar", pathCubeMx, "-s", projectFile)
		} else {
			cmd = exec.Command(pathJava, arg0, arg1, "-jar", pathCubeMx)
		}
	} else {
		if iocFile != "" {
			cmd = exec.Command(pathJava, "-jar", pathCubeMx, iocFile)
		} else if projectFile != "" {
			cmd = exec.Command(pathJava, "-jar", pathCubeMx, "-s", projectFile)
		} else {
			cmd = exec.Command(pathJava, "-jar", pathCubeMx)
		}
	}
	log.Debugf("Start CubeMX as %v", cmd)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return -1, err
	}

	return cmd.Process.Pid, nil
}

func WriteProjectFile(workDir string, params BridgeParamType) (string, error) {
	filePath := filepath.Join(workDir, "project.script")
	log.Debugf("Writing CubeMX project file %v", filePath)

	var text utils.TextBuilder
	if params.BoardName != "" && params.BoardVendor == "STMicroelectronics" {
		text.AddLine("loadboard", params.BoardName, "allmodes")
	} else {
		text.AddLine("load", params.Device)
	}
	text.AddLine("project name", "STM32CubeMX")

	toolchain, err := GetToolchain(params.Compiler)
	if err != nil {
		return "", err
	}
	text.AddLine("project toolchain", utils.AddQuotes(toolchain))

	cubeWorkDir := workDir
	if runtime.GOOS == "windows" {
		cubeWorkDir = filepath.FromSlash(cubeWorkDir)
	}
	text.AddLine("project path", utils.AddQuotes(cubeWorkDir))
	text.AddLine("SetCopyLibrary", utils.AddQuotes("copy only"))

	if utils.FileExists(filePath) {
		os.Remove(filePath)
	}

	err = os.WriteFile(filePath, []byte(text.GetLine()), 0600)
	if err != nil {
		log.Errorf("Error writing %v", err)
		return "", err
	}

	return filePath, nil
}

func ReadCbuildGenIdxYmlFile(path, generatorID string, parms *cbuild.ParamsType) error {
	log.Debugf("Reading cbuild-gen-idx.yml file: '%v'", path)
	err := cbuild.Read(path, generatorID, parms)
	if err != nil {
		return err
	}

	return nil
}

func ReadGeneratorYmlFile(path string, parms *generator.ParamsType) error {
	log.Debugf("Reading generator.yml file: '%v'", path)
	err := generator.Read(path, parms)
	return err
}

func GetBridgeInfo(parms *cbuild.ParamsType, bridgeParams *[]BridgeParamType) error {
	var boardName string
	var boardVendor string
	var device string
	var output string
	var projectType string

	split := strings.Split(parms.Board, "::")
	if len(split) == 2 {
		boardVendor = split[0]
		boardName = split[1]
	} else {
		boardVendor = ""
		boardName = parms.Board
	}
	split = strings.Split(boardName, ":")
	if len(split) == 2 {
		boardName = split[0]
	}

	device = parms.Device
	projectType = parms.ProjectType
	output = parms.Output

	for _, gen := range parms.CbuildGens {
		var bparm BridgeParamType

		bparm.BoardName = boardName
		bparm.BoardVendor = boardVendor
		bparm.Device = device
		bparm.Output = output
		bparm.ProjectName = gen.Project
		bparm.ProjectType = projectType
		bparm.ForProjectPart = gen.ForProjectPart
		bparm.GeneratorMap = gen.Map
		bparm.CgenName = gen.Name
		compiler := gen.CbuildGen.BuildGen.Compiler
		compiler = strings.Split(compiler, "@")[0]
		bparm.Compiler = compiler
		if gen.Map != "" {
			bparm.CubeContext = gen.Map
			bparm.CubeContextFolder = gen.Map
			if (parms.ProjectType == "trustzone") && (gen.ForProjectPart == "non-secure") {
				nonSecureSecurePairs := map[string]string{
					// can be extended
					"AppliNonSecure": "AppliSecure",
				}
				for _, tmpGen := range parms.CbuildGens {
					if (tmpGen.ForProjectPart == "secure") && (tmpGen.Map == nonSecureSecurePairs[gen.Map]) {
						bparm.PairedSecurePart = tmpGen.Project
					}
				}
			}
		} else {
			switch parms.ProjectType {
			case "single-core":
				bparm.CubeContext = ""
				bparm.CubeContextFolder = ""
			case "multi-core":
				core := gen.CbuildGen.BuildGen.Processor.Core

				cubeContext := strings.ReplaceAll(core, "-", "")
				cubeContext = strings.ReplaceAll(cubeContext, "+", "Plus") // Cortex-M0+ -> CortexM0Plus
				bparm.CubeContext = cubeContext

				cubeContextFolder := "C" + strings.Split(core, "-")[1]
				cubeContextFolder = strings.ReplaceAll(cubeContextFolder, "+", "PLUS") // Cortex-M0+ -> CM0PLUS
				bparm.CubeContextFolder = cubeContextFolder

			case "trustzone":
				core := gen.CbuildGen.BuildGen.Processor.Core
				context := strings.ReplaceAll(core, "-", "")
				if gen.ForProjectPart == "non-secure" {
					bparm.CubeContext = context + "NS"
					bparm.CubeContextFolder = "NonSecure"
					for _, tmpGen := range parms.CbuildGens {
						if tmpGen.ForProjectPart == "secure" {
							bparm.PairedSecurePart = tmpGen.Project
							break
						}
					}
				}
				if gen.ForProjectPart == "secure" {
					bparm.CubeContext = context + "S"
					bparm.CubeContextFolder = "Secure"
				}
			}
		}
		*bridgeParams = append(*bridgeParams, bparm)
	}
	return nil
}

var filterFiles = map[string]string{
	"system_":                            "system_ file (already added)",
	"Templates":                          "Templates file (mostly not present)",
	"/STM32CubeMX/Drivers/CMSIS/Include": "CMSIS include folder (delivered by ARM::CMSIS)",
}

func FilterFile(file string) bool {
	for key, value := range filterFiles {
		if strings.Contains(file, key) {
			log.Debugf("ignoring %v: %v", value, file)
			return true
		}
	}

	return false
}

func FindMxProject(context string, mxprojectAll MxprojectAllType) (MxprojectType, error) {
	if len(mxprojectAll.Mxproject) == 0 {
		return MxprojectType{}, errors.New("no .mxproject read")
	} else if len(mxprojectAll.Mxproject) == 1 {
		mxproject := mxprojectAll.Mxproject[0]
		return mxproject, nil
	}

	for _, mxproject := range mxprojectAll.Mxproject {
		if mxproject.Context == context {
			return mxproject, nil
		}
	}

	return MxprojectType{}, nil
}

func WriteCgenYml(outPath string, mxprojectAll MxprojectAllType, bridgeParams []BridgeParamType) error {
	for _, parm := range bridgeParams {
		mxproject, err := FindMxProject(parm.CubeContext, mxprojectAll)
		if err != nil {
			continue
		}
		err = WriteCgenYmlSub(outPath, mxproject, parm)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteCgenYmlSub(outPath string, mxproject MxprojectType, bridgeParam BridgeParamType) error {
	var cgen cbuild.CgenType

	relativePathAdd, err := GetRelativePathAdd(outPath, bridgeParam.Compiler)
	if err != nil {
		return err
	}

	cgen.GeneratorImport.ForBoard = bridgeParam.BoardName
	cgen.GeneratorImport.ForDevice = bridgeParam.Device
	cgen.GeneratorImport.Define = append(cgen.GeneratorImport.Define, mxproject.PreviousUsedFiles.CDefines...)

	for _, headerPath := range mxproject.PreviousUsedFiles.HeaderPath {
		headerPath, _ = utils.ConvertFilename(outPath, headerPath, relativePathAdd)
		if FilterFile(headerPath) {
			continue
		}
		cgen.GeneratorImport.AddPath = append(cgen.GeneratorImport.AddPath, headerPath)
	}

	cfgPath := "MX_Device"
	if bridgeParam.CubeContextFolder != "" {
		cfgPath = path.Join(cfgPath, bridgeParam.CubeContextFolder)
	}
	cfgPath, _ = utils.ConvertFilename(outPath, cfgPath, "")
	cgen.GeneratorImport.AddPath = append(cgen.GeneratorImport.AddPath, cfgPath)

	var groupSrc cbuild.CgenGroupsType
	var groupHalDriver cbuild.CgenGroupsType
	var groupTz cbuild.CgenGroupsType

	groupSrc.Group = "CubeMX"
	groupHalDriver.Group = "STM32 HAL Driver"
	groupHalFilter := "HAL_Driver"

	for _, file := range mxproject.PreviousUsedFiles.SourceFiles {
		if FilterFile(file) {
			continue
		}
		file, _ = utils.ConvertFilename(outPath, file, relativePathAdd)
		if strings.Contains(file, groupHalFilter) {
			var cgenFile cbuild.CgenFilesType
			cgenFile.File = file
			groupHalDriver.Files = append(groupHalDriver.Files, cgenFile)
		} else {
			var cgenFile cbuild.CgenFilesType
			cgenFile.File = file
			groupSrc.Files = append(groupSrc.Files, cgenFile)
		}
	}

	var cgenFile cbuild.CgenFilesType
	startupFile, err := GetStartupFile(outPath, bridgeParam)
	if err != nil {
		return err
	}
	startupFile, err = utils.ConvertFilenameRel(outPath, startupFile)
	if err != nil {
		return err
	}
	cgenFile.File = startupFile
	groupSrc.Files = append(groupSrc.Files, cgenFile)

	systemFile, err := GetSystemFile(outPath, bridgeParam)
	if err != nil {
		return err
	}
	systemFile, err = utils.ConvertFilenameRel(outPath, systemFile)
	if err != nil {
		return err
	}
	cgenFile.File = systemFile
	groupSrc.Files = append(groupSrc.Files, cgenFile)

	// linkerFiles, err := GetLinkerScripts(outPath, bridgeParam)
	// if err != nil {
	// 	return err
	// }
	// for _, file := range linkerFiles {
	// 	file, err = utils.ConvertFilenameRel(outPath, file)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	var cgenFile cbuild.CgenFilesType
	// 	cgenFile.File = file
	// 	groupSrc.Files = append(groupSrc.Files, cgenFile)
	// }

	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupSrc)
	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupHalDriver)

	if bridgeParam.ForProjectPart == "non-secure" {
		groupTz.Group = "CMSE Library"
		var cgenFile cbuild.CgenFilesType
		cgenFile.File = "$cmse-lib("
		cgenFile.File += bridgeParam.PairedSecurePart
		cgenFile.File += ")$"
		groupTz.Files = append(groupTz.Files, cgenFile)
		cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupTz)
	}

	return common.WriteYml(bridgeParam.CgenName, &cgen)
}

func GetToolchain(compiler string) (string, error) {
	var toolchainMapping = map[string]string{
		"AC6":   "MDK-ARM V5",
		"GCC":   "STM32CubeIDE",
		"IAR":   "EWARM",
		"CLANG": "STM32CubeIDE",
	}

	toolchain, ok := toolchainMapping[compiler]
	if !ok {
		return "", errors.New("unknown compiler '" + compiler + "'")
	}
	return toolchain, nil
}

func GetRelativePathAdd(outPath string, compiler string) (string, error) {
	var pathMapping = map[string]string{
		"AC6":   "MDK-ARM",
		"GCC":   "",
		"IAR":   "EWARM",
		"CLANG": "",
	}

	folder, ok := pathMapping[compiler]
	if !ok {
		return "", errors.New("unknown compiler '" + compiler + "'")
	}

	lastPath := filepath.Base(outPath)
	var relativePathAdd string
	if lastPath != "STM32CubeMX" {
		relativePathAdd = path.Join(relativePathAdd, "STM32CubeMX")
	}
	relativePathAdd = path.Join(relativePathAdd, folder)

	return relativePathAdd, nil
}

func GetToolchainFolderPath(outPath string, compiler string) (string, error) {
	var toolchainFolderMapping = map[string]string{
		"AC6":   "MDK-ARM",
		"GCC":   "STM32CubeIDE",
		"IAR":   "EWARM",
		"CLANG": "STM32CubeIDE",
	}

	toolchainFolder, ok := toolchainFolderMapping[compiler]
	if !ok {
		return "", errors.New("unknown compiler '" + compiler + "'")
	}

	lastPath := filepath.Base(outPath)
	toolchainFolderPath := outPath
	if lastPath != "STM32CubeMX" {
		toolchainFolderPath = path.Join(outPath, "STM32CubeMX")
	}
	toolchainFolderPath = path.Join(toolchainFolderPath, toolchainFolder)

	return toolchainFolderPath, nil
}

func GetStartupFile(outPath string, bridgeParams BridgeParamType) (string, error) {
	var startupFolder string
	var fileFilter string
	var fileExtensions = []string{".s", ".S", ".c"}

	startupFolder, err := GetToolchainFolderPath(outPath, bridgeParams.Compiler)
	if err != nil {
		return "", err
	}

	if bridgeParams.CubeContextFolder != "" {
		fileFilter = "_" + bridgeParams.CubeContextFolder
	}

	if bridgeParams.Compiler == "GCC" || bridgeParams.Compiler == "CLANG" {
		startupFolder = path.Join(startupFolder, bridgeParams.CubeContextFolder, "Application", "Startup")
	}

	if !utils.DirExists(startupFolder) {
		errorString := "Directory not found: " + startupFolder
		log.Error(errorString)
		return "", errors.New(errorString)
	}

	var startupFile string
	var defaultStartupFile string
	var startupFileList []string

	err = filepath.Walk(startupFolder, func(path string, f fs.FileInfo, err error) error {
		if f.Mode().IsRegular() &&
			strings.HasPrefix(f.Name(), "startup_") &&
			(strings.HasSuffix(f.Name(), fileExtensions[0]) ||
				strings.HasSuffix(f.Name(), fileExtensions[1]) ||
				strings.HasSuffix(f.Name(), fileExtensions[2])) {

			startupFileList = append(startupFileList, path)
		}
		return nil
	})

	if len(startupFileList) == 1 {
		startupFile = startupFileList[0]
	} else if len(startupFileList) > 1 {
		for _, file := range startupFileList {
			fileName := strings.Split(filepath.Base(file), ".")[0]
			split := strings.Split(fileName, "_")
			if len(split) == 2 {
				defaultStartupFile = file
			} else {
				fileFilterLower := strings.ToLower(fileFilter)
				nameLower := strings.ToLower(fileName)
				if strings.Contains(nameLower, fileFilterLower) {
					startupFile = file
					break
				}
			}
		}
		if startupFile == "" {
			startupFile = defaultStartupFile
		}
	}

	if startupFile == "" {
		errorString := "startup file not found"
		log.Error(errorString)
		return "", errors.New(errorString)
	}

	return startupFile, err
}

func GetSystemFile(outPath string, bridgeParams BridgeParamType) (string, error) {
	var toolchainFolder string
	var systemFolder string

	toolchainFolder, err := GetToolchainFolderPath(outPath, bridgeParams.Compiler)
	if err != nil {
		return "", err
	}

	if bridgeParams.ProjectType == "multi-core" {
		systemFolder = filepath.Dir(toolchainFolder)
		systemFolder = path.Join(systemFolder, "Common")
		if !utils.DirExists(toolchainFolder) {
			systemFolder = ""
		}
	}

	if systemFolder == "" {
		systemFolder = filepath.Dir(toolchainFolder)

		if bridgeParams.CubeContextFolder != "" {
			systemFolder = path.Join(systemFolder, bridgeParams.CubeContextFolder)
		}

		systemFolder = path.Join(systemFolder, "Src")
	}

	if !utils.DirExists(systemFolder) {
		errorString := "Directory not found: " + systemFolder
		log.Error(errorString)
		return "", errors.New(errorString)
	}

	var systemFile string
	err = filepath.Walk(systemFolder, func(path string, f fs.FileInfo, err error) error {
		if f.Mode().IsRegular() &&
			strings.HasPrefix(f.Name(), "system_stm32") &&
			strings.HasSuffix(f.Name(), ".c") {
			systemFile = path
		}
		return nil
	})

	if systemFile == "" {
		errorString := "system file not found"
		log.Error(errorString)
		return "", errors.New(errorString)
	}

	return systemFile, err
}

// func GetLinkerScripts(outPath string, bridgeParams BridgeParamType) ([]string, error) {
// 	var linkerFolder string
// 	var fileExtesion string
// 	var fileFilter string

// 	linkerFolder, err := GetToolchainFolderPath(outPath, bridgeParams.Compiler)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch bridgeParams.Compiler {
// 	case "AC6":
// 		fileExtesion = ".sct"
// 	case "IAR":
// 		fileExtesion = ".icf"
// 	case "GCC", "CLANG":
// 		fileExtesion = ".ld"
// 	default:
// 		return nil, errors.New("unknown compiler '" + bridgeParams.Compiler + "'")
// 	}

// 	if bridgeParams.GeneratorMap != "" {
// 		linkerFolder = path.Join(linkerFolder, bridgeParams.GeneratorMap)
// 	}

// 	switch bridgeParams.Compiler {
// 	case "AC6", "IAR":
// 		switch bridgeParams.ProjectType {
// 		case "single-core":
// 			fileFilter = ""
// 		case "multi-core":
// 			fileFilter = "_" + bridgeParams.ForProjectPart
// 		case "trustzone":
// 			if bridgeParams.ForProjectPart == "secure" {
// 				fileFilter = "_s."
// 			}
// 			if bridgeParams.ForProjectPart == "non-secure" {
// 				fileFilter = "_ns."
// 			}
// 		}

// 	case "GCC", "CLANG":
// 		switch bridgeParams.ProjectType {
// 		case "multi-core":
// 			linkerFolder = path.Join(linkerFolder, bridgeParams.ForProjectPart)
// 		case "trustzone":
// 			if bridgeParams.ForProjectPart == "secure" {
// 				linkerFolder = path.Join(linkerFolder, "Secure")
// 			}
// 			if bridgeParams.ForProjectPart == "non-secure" {
// 				linkerFolder = path.Join(linkerFolder, "NonSecure")
// 			}
// 		}
// 	default:
// 		return nil, errors.New("unknown compiler '" + bridgeParams.Compiler + "'")
// 	}

// 	if !utils.DirExists(linkerFolder) {
// 		errorString := "Directory not found: " + linkerFolder
// 		log.Error(errorString)
// 		return nil, errors.New(errorString)
// 	}

// 	var linkerScripts []string
// 	err = filepath.Walk(linkerFolder, func(path string, f fs.FileInfo, err error) error {
// 		if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), fileExtesion) {
// 			if fileFilter != "" {
// 				if strings.Contains(f.Name(), fileFilter) {
// 					linkerScripts = append(linkerScripts, path)
// 				}
// 			} else {
// 				linkerScripts = append(linkerScripts, path)
// 			}
// 		}
// 		return nil
// 	})

// 	return linkerScripts, err
// }
