/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmlPath, outPath, cubeMxPath string) error {
	var projectFile string
	var parms cbuild.ParamsType

	err := ReadCbuildYmlFile(cbuildYmlPath, outPath, &parms)
	if err != nil {
		return err
	}

	workDir := path.Dir(cbuildYmlPath)
	if parms.OutPath != "" {
		workDir = path.Join(workDir, parms.OutPath)
	} else {
		workDir = path.Join(workDir, outPath)
	}
	workDir = filepath.Clean(workDir)
	workDir = filepath.ToSlash(workDir)

	err = os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		return err
	}

	cubeIocPath := path.Join(workDir, "STM32CubeMX", "STM32CubeMX.ioc")

	if utils.FileExists(cubeIocPath) {
		err := Launch(cubeIocPath, "")
		if err != nil {
			return err
		}
	} else {
		projectFile, err = WriteProjectFile(workDir, &parms)
		if err != nil {
			return nil
		}
		log.Infof("Generated file: %v", projectFile)

		err := Launch("", projectFile)
		if err != nil {
			return err
		}
	}

	mxprojectPath := path.Join(workDir, "STM32CubeMX", ".mxproject")
	mxproject, err := IniReader(mxprojectPath, false)
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
	text.AddLine("project toolchain", utils.AddQuotes("MDK-ARM V5"))
	cubeWorkDir := workDir
	if runtime.GOOS == "windows" {
		cubeWorkDir = filepath.FromSlash(cubeWorkDir)
	}
	text.AddLine("project path", utils.AddQuotes(cubeWorkDir))
	text.AddLine("SetCopyLibrary", utils.AddQuotes("copy only"))

	if utils.FileExists(filePath) {
		os.Remove(filePath)
	}

	err := os.WriteFile(filePath, []byte(text.GetLine()), 0600)
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

func ConvertFilename(outPath, file string) (string, error) {
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	mdkarmPath := path.Join(outPath, "STM32CubeMX", "MDK-ARM") // create the path where STCube sets it's files relative to (STM32CubeMX/MDK-ARM/)
	file = path.Join(mdkarmPath, file)

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Errorf("file or directory not found: %v", file)
	}

	var err error
	origfilename := file
	file, err = filepath.Rel(outPath, file)
	if err != nil {
		log.Errorf("path error found: %v", file)
		return origfilename, errors.New("path error")
	}

	file = filepath.ToSlash(file)
	file = "./" + file

	return file, nil
}

func WriteCgenYml(outPath string, mxproject MxprojectType, inParms cbuild.ParamsType) error {
	outFile := path.Join(outPath, "STM32CubeMX.cgen.yml")
	var cgen cbuild.CgenType

	cgen.Layer.ForBoard = inParms.Board
	cgen.Layer.ForDevice = inParms.Device
	cgen.Layer.Define = append(cgen.Layer.Define, mxproject.PreviousUsedKeilFiles.CDefines...)

	for id := range mxproject.PreviousUsedKeilFiles.HeaderPath {
		headerPath := mxproject.PreviousUsedKeilFiles.HeaderPath[id]
		headerPath, _ = ConvertFilename(outPath, headerPath)
		cgen.Layer.AddPath = append(cgen.Layer.AddPath, headerPath)
	}

	var groupSrc cbuild.CgenGroupsType
	var groupHalDriver cbuild.CgenGroupsType
	groupSrc.Group = "CubeMX"
	groupHalDriver.Group = "HAL Driver"
	groupHalFilter := "HAL_Driver"

	for id := range inParms.Subsystem {
		subsystem := inParms.Subsystem[id]
		packs := subsystem.Packs

		for id2 := range packs {
			pack := packs[id2]
			var cgenPack cbuild.CgenPacksType
			cgenPack.Pack = pack.Pack
			cgen.Layer.Packs = append(cgen.Layer.Packs, cgenPack)
		}
	}

	for id := range mxproject.PreviousUsedKeilFiles.SourceFiles {
		file := mxproject.PreviousUsedKeilFiles.SourceFiles[id]
		file, _ = ConvertFilename(outPath, file)

		if strings.Contains(file, "Templates") {
			log.Infof("ignoring Templates file (mostly not present): %v", file)
			continue
		}

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

	cgen.Layer.Groups = append(cgen.Layer.Groups, groupSrc)
	cgen.Layer.Groups = append(cgen.Layer.Groups, groupHalDriver)

	return common.WriteYml(outFile, &cgen)
}
