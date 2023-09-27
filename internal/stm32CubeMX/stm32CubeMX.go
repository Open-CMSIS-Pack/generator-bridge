/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmlPath, outPath, cubeMxPath string) error {
	var projectFile string
	var parms cbuild.ParamsType

	ReadCbuildYmlFile(cbuildYmlPath, outPath, &parms)

	workDir := path.Dir(cbuildYmlPath)
	if parms.OutPath != "" {
		workDir = path.Join(workDir, parms.OutPath)
	} else {
		workDir = path.Join(workDir, outPath)
	}
	workDir = filepath.Clean(workDir)
	workDir = filepath.ToSlash(workDir)

	err := os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		return err
	}

	cubeIocPath := path.Join(workDir, "STM32CubeMX", "STM32CubeMX.ioc")

	if utils.FileExists(cubeIocPath) {
		Launch(cubeIocPath, "")
	} else {
		projectFile, err = WriteProjectFile(workDir, &parms)
		if err != nil {
			return nil
		}
		log.Infof("Generated file: %v", projectFile)

		Launch("", projectFile)
	}

	mxprojectPath := path.Join(workDir, "STM32CubeMX", ".mxproject")
	mxproject, err := IniReader(mxprojectPath, false)
	if err != nil {
		return err
	}

	WriteCgenYml(workDir, mxproject, parms)

	return nil
}

func Launch(iocFile, projectFile string) error {
	log.Infof("Launching STM32CubeMX...")

	pathJava := path.Join(os.Getenv("STM32CubeMX_PATH"), "jre", "bin", "java.exe")
	pathCubeMx := path.Join(os.Getenv("STM32CubeMX_PATH"), "STM32CubeMX.exe")

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
	text.AddLine("project path", utils.AddQuotes(workDir))
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

func WriteCgenYml(outPath string, mxproject MxprojectType, inParms cbuild.ParamsType) error {
	outFile := path.Join(outPath, "STM32CubeMX.cgen.yml")
	var cgen cbuild.CgenType

	cgen.Layer.ForBoard = inParms.Board
	cgen.Layer.ForDevice = inParms.Device
	cgen.Layer.Define = append(cgen.Layer.Define, mxproject.PreviousUsedKeilFiles.CDefines...)
	cgen.Layer.AddPath = append(cgen.Layer.AddPath, mxproject.PreviousUsedKeilFiles.HeaderPath...)

	var groupSrc cbuild.CgenGroupsType
	var groupHalDriver cbuild.CgenGroupsType
	groupSrc.Group = "STM32CubeMX"
	groupHalDriver.Group = "HAL_Driver"

	for id := range inParms.Core {
		core := inParms.Core[id]
		packs := core.Packs

		for id2 := range packs {
			pack := packs[id2]
			var cgenPack cbuild.CgenPacksType
			cgenPack.Pack = pack.Pack
			cgen.Layer.Packs = append(cgen.Layer.Packs, cgenPack)
		}
	}

	for id := range mxproject.PreviousUsedKeilFiles.SourceFiles {
		file := mxproject.PreviousUsedKeilFiles.SourceFiles[id]
		if strings.Contains(file, "HAL_Driver") {
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

	var header utils.TextBuilder
	//header.AddLine("# File Name   :", filepath.Base(outFile))
	//header.AddLine("# Date        :", utils.GetDateTimeString())
	//header.AddLine("# Description :", "Generator layer")
	//header.AddLine("#")
	//header.AddLine("")

	return common.WriteYml(outFile, header.GetLine(), &cgen)
}
