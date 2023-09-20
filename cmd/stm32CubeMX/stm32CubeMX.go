/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package stm32CubeMX

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/open-cmsis-pack/generator-bridge/cmd/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/cmd/utils"
	log "github.com/sirupsen/logrus"
)

func Process(cbuildYmpPath, cubeMxPath string) error {
	var projectFile string
	var err error

	cubeIocPath := path.Join(path.Dir(cbuildYmpPath), "STM32CubeMX", "STM32CubeMX.ioc")
	if utils.FileExists(cubeIocPath) {
		//Launch(cubeMxPath, cubeIocPath, "")
	} else {
		var parms cbuild.Params_s
		ReadCbuildYmlFile(cbuildYmpPath, &parms)
		workDir := path.Dir(cbuildYmpPath)
		projectFile, err = WriteProjectFile(workDir, &parms)
		if err != nil {
			return nil
		}

		Launch(cubeMxPath, "", projectFile)
	}

	mxprojectPath := path.Join(path.Dir(cbuildYmpPath), "STM32CubeMX", ".mxproject")
	IniReader(mxprojectPath)

	return nil
}

func Launch(pathCubeMx, iocFile, projectFile string) error {
	log.Infof("Launching STM32CubeMX...")

	if pathCubeMx == "" {
		pathCubeMx = path.Join(os.Getenv("STM32CubeMX_PATH"), "STM32CubeMX.exe")
	}

	var cmd *exec.Cmd
	if iocFile != "" {
		cmd = exec.Command(pathCubeMx, iocFile)
	} else if projectFile != "" {
		cmd = exec.Command(pathCubeMx, "-s", projectFile)
	} else {
		cmd = exec.Command(pathCubeMx)
	}
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Press ENTER to continue...")
	//reader := bufio.NewReader(os.Stdin)
	//reader.ReadString('\n')
	os.Stdin.Read([]byte{0})

	return nil
}

func WriteProjectFile(workDir string, parms *cbuild.Params_s) (string, error) {
	filePath := filepath.Join(workDir, "project.script")
	log.Infof("Writing CubeMX project file %v", filePath)

	var text utils.TextBuilder
	text.AddLine("loadboard", parms.Board, "allmodes")
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
