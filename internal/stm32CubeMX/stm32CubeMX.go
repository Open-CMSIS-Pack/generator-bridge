/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/generator"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmlPath, outPath, cubeMxPath, mxprojectPath string, runCubeMx bool) error {
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
	err = ReadCbuildYmlFile(cbuildYmlPath, outPath, &parms)
	if err != nil {
		return err
	}

	workDir := path.Dir(cbuildYmlPath)
	if parms.OutPath != "" {
		if filepath.IsAbs(parms.OutPath) {
			workDir = parms.OutPath
		} else {
			workDir = path.Join(workDir, parms.OutPath)
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

	if runCubeMx {
		cubeIocPath := workDir
		lastPath := filepath.Base(cubeIocPath)
		if lastPath != "STM32CubeMX" {
			cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX")
		}
		cubeIocPath = path.Join(cubeIocPath, "STM32CubeMX.ioc")

		if utils.FileExists(cubeIocPath) {
			err := Launch(cubeIocPath, "")
			if err != nil {
				return errors.New("generator '" + gParms.ID + "' missing. Install from '" + gParms.DownloadURL + "'")
			}
		} else {
			projectFile, err = WriteProjectFile(workDir, &parms)
			if err != nil {
				return nil
			}
			log.Infof("Generated file: %v", projectFile)

			err := Launch("", projectFile)
			if err != nil {
				return errors.New("generator '" + gParms.ID + "' missing. Install from '" + gParms.DownloadURL + "'")
			}
		}

		err = ReadContexts(cubeIocPath, parms)
		if err != nil {
			return err
		}

		tmpPath, _ := filepath.Split(cubeIocPath)
		mxprojectPath = path.Join(tmpPath, ".mxproject")
	}
	mxproject, err := IniReader(mxprojectPath, parms.Subsystem[0].Compiler, false)
	if err != nil {
		return err
	}

	err = WriteCgenYml(workDir, mxproject, parms)
	if err != nil {
		return err
	}

	return nil
}

func Launch(iocFile, projectFile string) error {
	log.Infof("Launching STM32CubeMX...")

	const cubeEnvVar = "STM32CubeMX_PATH"
	cubeEnv := os.Getenv(cubeEnvVar)
	if cubeEnv == "" {
		return errors.New("environment variable for CubeMX not set: " + cubeEnvVar)
	}

	pathJava := path.Join(cubeEnv, "jre", "bin", "java.exe")
	pathCubeMx := path.Join(cubeEnv, "STM32CubeMX.exe")

	var cmd *exec.Cmd
	if iocFile != "" {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx, iocFile)
	} else if projectFile != "" {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx, "-s", projectFile)
	} else {
		cmd = exec.Command(pathJava, "-jar", pathCubeMx)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func WriteProjectFile(workDir string, parms *cbuild.ParamsType) (string, error) {
	filePath := filepath.Join(workDir, "project.script")
	log.Infof("Writing CubeMX project file %v", filePath)

	var text utils.TextBuilder
	if parms.Board != "" {
		text.AddLine("loadboard", parms.Board, "allmodes")
	} else {
		text.AddLine("load", parms.Device)
	}
	text.AddLine("project name", "STM32CubeMX")

	toolchain, err := GetToolchain(parms.Subsystem[0].Compiler)
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
		return "", err
	}

	return filePath, nil
}

func ReadCbuildYmlFile(path, outPath string, parms *cbuild.ParamsType) error {
	log.Infof("Reading cbuild.yml file: '%v'", path)
	err := cbuild.Read(path, outPath, parms)
	if err != nil {
		return err
	}

	return nil
}

func ReadGeneratorYmlFile(path string, parms *generator.ParamsType) error {
	log.Infof("Reading generator.yml file: '%v'", path)
	err := generator.Read(path, parms)
	return err
}

var filterFiles = map[string]string{
	"system_":                            "system_ file (already added)",
	"Templates":                          "Templates file (mostly not present)",
	"/STM32CubeMX/Drivers/CMSIS/Include": "CMSIS include folder (delivered by ARM::CMSIS)",
}

func FilterFile(file string) bool {
	for key, value := range filterFiles {
		if strings.Contains(file, key) {
			log.Infof("ignoring %v: %v", value, file)
			return true
		}
	}

	return false
}

func FindMxProject(subsystem *cbuild.SubsystemType, mxprojectAll MxprojectAllType) (MxprojectType, error) {
	if len(mxprojectAll.Mxproject) == 0 {
		return MxprojectType{}, errors.New("no .mxproject read")
	} else if len(mxprojectAll.Mxproject) == 1 {
		mxproject := mxprojectAll.Mxproject[0]
		return mxproject, nil
	}

	coreName := subsystem.CoreName
	trustzone := subsystem.TrustZone
	if trustzone == "off" {
		trustzone = ""
	}
	for id := range mxprojectAll.Mxproject {
		mxproject := mxprojectAll.Mxproject[id]
		if mxproject.CoreName == coreName && mxproject.Trustzone == trustzone {
			return mxproject, nil
		}
	}

	return MxprojectType{}, nil
}

func WriteCgenYml(outPath string, mxprojectAll MxprojectAllType, inParms cbuild.ParamsType) error {
	for id := range inParms.Subsystem {
		subsystem := &inParms.Subsystem[id]

		mxproject, err := FindMxProject(subsystem, mxprojectAll)
		if err != nil {
			continue
		}
		err = WriteCgenYmlSub(outPath, mxproject, subsystem)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteCgenYmlSub(outPath string, mxproject MxprojectType, subsystem *cbuild.SubsystemType) error {
	outName := subsystem.SubsystemIdx.Project + ".cgen.yml"
	outFile := path.Join(outPath, outName)
	var cgen cbuild.CgenType

	relativePathAdd, err := GetRelativePathAdd(outPath, subsystem.Compiler)
	if err != nil {
		return err
	}

	cgen.GeneratorImport.ForBoard = subsystem.Board
	cgen.GeneratorImport.ForDevice = subsystem.Device
	cgen.GeneratorImport.Define = append(cgen.GeneratorImport.Define, mxproject.PreviousUsedFiles.CDefines...)

	for _, headerPath := range mxproject.PreviousUsedFiles.HeaderPath {
		headerPath, _ = utils.ConvertFilename(outPath, headerPath, relativePathAdd)
		if FilterFile(headerPath) {
			continue
		}
		cgen.GeneratorImport.AddPath = append(cgen.GeneratorImport.AddPath, headerPath)
	}

	cfgPath := path.Join("drv_cfg", subsystem.SubsystemIdx.Project)
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
	startupFile, err := GetStartupFile(outPath, subsystem)
	if err != nil {
		return err
	}
	startupFile, err = utils.ConvertFilenameRel(outPath, startupFile)
	if err != nil {
		return err
	}
	cgenFile.File = startupFile
	groupSrc.Files = append(groupSrc.Files, cgenFile)

	systemFile, err := GetSystemFile(outPath, subsystem)
	if err != nil {
		return err
	}
	systemFile, err = utils.ConvertFilenameRel(outPath, systemFile)
	if err != nil {
		return err
	}
	cgenFile.File = systemFile
	groupSrc.Files = append(groupSrc.Files, cgenFile)

	linkerFiles, err := GetLinkerScripts(outPath, subsystem)
	if err != nil {
		return err
	}
	for _, file := range linkerFiles {
		file, err = utils.ConvertFilenameRel(outPath, file)
		if err != nil {
			return err
		}
		var cgenFile cbuild.CgenFilesType
		cgenFile.File = file
		groupSrc.Files = append(groupSrc.Files, cgenFile)
	}

	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupSrc)
	cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupHalDriver)

	if subsystem.TrustZone == "non-secure" {
		groupTz.Group = "CMSE Library"
		var cgenFile cbuild.CgenFilesType
		cgenFile.File = "$cmse-lib("
		cgenFile.File += subsystem.SubsystemIdx.SecureContextName
		cgenFile.File += ")$"
		groupTz.Files = append(groupTz.Files, cgenFile)
		cgen.GeneratorImport.Groups = append(cgen.GeneratorImport.Groups, groupTz)
	}

	return common.WriteYml(outFile, &cgen)
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
		return "", errors.New("unknown compiler")
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
		return "", errors.New("unknown compiler")
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
		return "", errors.New("unknown compiler")
	}

	lastPath := filepath.Base(outPath)
	toolchainFolderPath := outPath
	if lastPath != "STM32CubeMX" {
		toolchainFolderPath = path.Join(outPath, "STM32CubeMX")
	}
	toolchainFolderPath = path.Join(toolchainFolderPath, toolchainFolder)

	return toolchainFolderPath, nil
}

func GetStartupFile(outPath string, subsystem *cbuild.SubsystemType) (string, error) {
	var startupFolder string
	var fileExtesion string
	var fileFilter string

	startupFolder, err := GetToolchainFolderPath(outPath, subsystem.Compiler)
	if err != nil {
		return "", err
	}

	fileExtesion = ".s"
	switch subsystem.Compiler {
	case "AC6", "IAR":
		if subsystem.SubsystemIdx.ProjectType == "multi-core" {
			fileFilter = "_" + subsystem.SubsystemIdx.ForProjectPart
		}

	case "GCC", "CLANG":
		switch subsystem.SubsystemIdx.ProjectType {
		case "multi-core":
			startupFolder = path.Join(startupFolder, subsystem.SubsystemIdx.ForProjectPart)
		case "trustzone":
			if subsystem.SubsystemIdx.ForProjectPart == "secure" {
				startupFolder = path.Join(startupFolder, "Secure")
			}
			if subsystem.SubsystemIdx.ForProjectPart == "non-secure" {
				startupFolder = path.Join(startupFolder, "NonSecure")
			}
		}
		startupFolder = path.Join(startupFolder, "Application")
		startupFolder = path.Join(startupFolder, "Startup")

	default:
		return "", errors.New("unknown compiler")
	}

	if !utils.DirExists(startupFolder) {
		errorString := "Directory not found: " + startupFolder
		log.Errorf(errorString)
		return "", errors.New(errorString)
	}

	var startupFile string
	err = filepath.Walk(startupFolder, func(path string, f fs.FileInfo, err error) error {
		if f.Mode().IsRegular() &&
			strings.HasSuffix(f.Name(), fileExtesion) &&
			strings.HasPrefix(f.Name(), "startup_") {
			if fileFilter != "" {
				if strings.Contains(f.Name(), fileFilter) {
					startupFile = path
				}
			} else {
				startupFile = path
			}
		}
		return nil
	})

	if startupFile == "" {
		errorString := "startup file not found"
		log.Errorf(errorString)
		return "", errors.New(errorString)
	}

	return startupFile, err
}

func GetSystemFile(outPath string, subsystem *cbuild.SubsystemType) (string, error) {
	var toolchainFolder string
	var systemFolder string

	toolchainFolder, err := GetToolchainFolderPath(outPath, subsystem.Compiler)
	if err != nil {
		return "", err
	}

	if subsystem.SubsystemIdx.ProjectType == "multi-core" {
		systemFolder = filepath.Dir(toolchainFolder)
		systemFolder = path.Join(systemFolder, "Common")
		systemFolder = path.Join(systemFolder, "Src")
		if !utils.DirExists(toolchainFolder) {
			systemFolder = ""
		}
	}

	if systemFolder == "" {
		systemFolder = filepath.Dir(toolchainFolder)
		switch subsystem.SubsystemIdx.ProjectType {
		case "multi-core":
			systemFolder = path.Join(systemFolder, subsystem.SubsystemIdx.ForProjectPart)
		case "trustzone":
			if subsystem.SubsystemIdx.ForProjectPart == "secure" {
				systemFolder = path.Join(systemFolder, "Secure")
			}
			if subsystem.SubsystemIdx.ForProjectPart == "non-secure" {
				systemFolder = path.Join(systemFolder, "NonSecure")
			}
		}
		systemFolder = path.Join(systemFolder, "Src")
	}

	if !utils.DirExists(systemFolder) {
		errorString := "Directory not found: " + systemFolder
		log.Errorf(errorString)
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
		log.Errorf(errorString)
		return "", errors.New(errorString)
	}

	return systemFile, err
}

func GetLinkerScripts(outPath string, subsystem *cbuild.SubsystemType) ([]string, error) {
	var linkerFolder string
	var fileExtesion string
	var fileFilter string

	linkerFolder, err := GetToolchainFolderPath(outPath, subsystem.Compiler)
	if err != nil {
		return nil, err
	}

	switch subsystem.Compiler {
	case "AC6":
		fileExtesion = ".sct"
	case "IAR":
		fileExtesion = ".icf"
	case "GCC", "CLANG":
		fileExtesion = ".ld"
	default:
		return nil, errors.New("unknown compiler")
	}

	switch subsystem.Compiler {
	case "AC6", "IAR":
		switch subsystem.SubsystemIdx.ProjectType {
		case "single-core":
			fileFilter = ""
		case "multi-core":
			fileFilter = "_" + subsystem.SubsystemIdx.ForProjectPart
		case "trustzone":
			if subsystem.SubsystemIdx.ForProjectPart == "secure" {
				fileFilter = "_s."
			}
			if subsystem.SubsystemIdx.ForProjectPart == "non-secure" {
				fileFilter = "_ns."
			}
		}

	case "GCC", "CLANG":
		switch subsystem.SubsystemIdx.ProjectType {
		case "multi-core":
			linkerFolder = path.Join(linkerFolder, subsystem.SubsystemIdx.ForProjectPart)
		case "trustzone":
			if subsystem.SubsystemIdx.ForProjectPart == "secure" {
				linkerFolder = path.Join(linkerFolder, "Secure")
			}
			if subsystem.SubsystemIdx.ForProjectPart == "non-secure" {
				linkerFolder = path.Join(linkerFolder, "NonSecure")
			}
		}
	default:
		return nil, errors.New("unknown compiler")
	}

	if !utils.DirExists(linkerFolder) {
		errorString := "Directory not found: " + linkerFolder
		log.Errorf(errorString)
		return nil, errors.New(errorString)
	}

	var linkerScripts []string
	err = filepath.Walk(linkerFolder, func(path string, f fs.FileInfo, err error) error {
		if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), fileExtesion) {
			if fileFilter != "" {
				if strings.Contains(f.Name(), fileFilter) {
					linkerScripts = append(linkerScripts, path)
				}
			} else {
				linkerScripts = append(linkerScripts, path)
			}
		}
		return nil
	})

	return linkerScripts, err
}
