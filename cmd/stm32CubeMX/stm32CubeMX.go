/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package stm32CubeMX

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/cmd/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/cmd/common"
	"github.com/open-cmsis-pack/generator-bridge/cmd/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmlPath, cubeMxPath string) error {
	var projectFile string
	var err error

	cubeIocPath := path.Join(path.Dir(cbuildYmlPath), "STM32CubeMX", "STM32CubeMX.ioc")
	if utils.FileExists(cubeIocPath) {
		Launch(cubeIocPath, "")
	} else {
		var parms cbuild.Params_s
		ReadCbuildYmlFile(cbuildYmlPath, &parms)
		workDir := path.Dir(cbuildYmlPath)
		projectFile, err = WriteProjectFile(workDir, &parms)
		if err != nil {
			return nil
		}

		Launch("", projectFile)
	}

	mxprojectPath := path.Join(path.Dir(cbuildYmlPath), "STM32CubeMX", ".mxproject")
	IniReader(mxprojectPath, false)

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

func ReadCbuildYmlFile(path string, parms *cbuild.Params_s) error {
	log.Infof("Reading cbuild.yml file: '%v'", path)
	cbuild.Read(path, parms)

	return nil
}

const header string = `#
# File Name   : STM32CubeMX.cgen.yml
# Date        : 29/08/2023 07:05:15
# Description : Generator layer
#
`

func WriteCgenYml(outPath string, mxproject Mxproject_s) error {
	outFile := path.Join(outPath, "STM32CubeMX.cgen.yml")
	var cgen cbuild.Cgen_s

	cgen.Generator.ForBoard = ""
	cgen.Generator.ForDevice = "STM32"
	cgen.Generator.GeneratedBy = "STM32CubeMX bridge"
	cgen.Generator.Define = append(cgen.Generator.Define, mxproject.PreviousUsedKeilFiles.CDefines...)
	cgen.Generator.AddPath = append(cgen.Generator.AddPath, mxproject.PreviousUsedKeilFiles.HeaderPath...)

	var groupSrc cbuild.CgenGroups_s
	var groupHalDriver cbuild.CgenGroups_s
	groupSrc.Group = "STM32CubeMX"
	groupHalDriver.Group = "HAL_Driver"

	for id := range mxproject.PreviousUsedKeilFiles.SourceFiles {
		file := mxproject.PreviousUsedKeilFiles.SourceFiles[id]
		if !strings.Contains(file, "HAL_Driver") {
			var cgenFile cbuild.CgenFiles_s
			cgenFile.File = file
			groupSrc.Files = append(groupSrc.Files, cgenFile)
		} else {
			var cgenFile cbuild.CgenFiles_s
			cgenFile.File = file
			groupHalDriver.Files = append(groupSrc.Files, cgenFile)
		}
	}

	cgen.Generator.Groups = append(cgen.Generator.Groups, groupSrc)
	cgen.Generator.Groups = append(cgen.Generator.Groups, groupHalDriver)

	return common.WriteYml(outFile, header, &cgen)
}
