/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package stm32CubeMX

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
	var parms cbuild.Params_s

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

func WriteProjectFile(workDir string, parms *cbuild.Params_s) (string, error) {
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

	os.WriteFile(filePath, []byte(text.GetLine()), 0777)

	return filePath, nil
}

func ReadCbuildYmlFile(path, outPath string, parms *cbuild.Params_s) error {
	log.Infof("Reading cbuild.yml file: '%v'", path)
	cbuild.Read(path, outPath, parms)

	return nil
}

func WriteCgenYml(outPath string, mxproject Mxproject_s, inParms cbuild.Params_s) error {
	outFile := path.Join(outPath, "STM32CubeMX.cgen.yml")
	var cgen cbuild.Cgen_s

	cgen.Layer.ForBoard = inParms.Board
	cgen.Layer.ForDevice = inParms.Device
	cgen.Layer.Define = append(cgen.Layer.Define, mxproject.PreviousUsedKeilFiles.CDefines...)
	cgen.Layer.AddPath = append(cgen.Layer.AddPath, mxproject.PreviousUsedKeilFiles.HeaderPath...)

	var groupSrc cbuild.CgenGroups_s
	var groupHalDriver cbuild.CgenGroups_s
	groupSrc.Group = "STM32CubeMX"
	groupHalDriver.Group = "HAL_Driver"

	for id := range inParms.Core {
		core := inParms.Core[id]
		packs := core.Packs

		for id2 := range packs {
			pack := packs[id2]
			var cgenPack cbuild.CgenPacks_s
			cgenPack.Pack = pack.Pack
			cgen.Layer.Packs = append(cgen.Layer.Packs, cgenPack)
		}
	}

	for id := range mxproject.PreviousUsedKeilFiles.SourceFiles {
		file := mxproject.PreviousUsedKeilFiles.SourceFiles[id]
		if strings.Contains(file, "HAL_Driver") {
			var cgenFile cbuild.CgenFiles_s
			cgenFile.File = file
			groupHalDriver.Files = append(groupHalDriver.Files, cgenFile)
		} else {
			var cgenFile cbuild.CgenFiles_s
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
