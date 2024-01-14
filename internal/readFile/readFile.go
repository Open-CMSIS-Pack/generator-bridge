/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package readfile

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	stm32cubemx "github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	log "github.com/sirupsen/logrus"
)

func Process(inFile, inFile2, outPath string) error {
	log.Infof("Reading file: %v", inFile)
	if outPath == "" {
		outPath = filepath.Dir(inFile)
	}

	var params cbuild.ParamsType

	if strings.Contains(inFile, "cbuild-gen-idx.yml") {
		err := cbuild.Read(inFile, outPath, &params)
		if err != nil {
			return err
		}
		_, err = stm32cubemx.WriteProjectFile(outPath, &params)
		if err != nil {
			return err
		}
	} else if strings.Contains(inFile, "cbuild-gen.yml") || strings.Contains(inFile, "cbuild.yml") {
		params.OutPath = outPath
		var subsystem cbuild.SubsystemType
		err := cbuild.ReadCbuildgen(inFile, &subsystem)
		if err != nil {
			return err
		}
		params.Subsystem = append(params.Subsystem, subsystem)
	}

	var mxprojectFile string
	if strings.Contains(inFile, ".mxproject") {
		mxprojectFile = inFile
	}
	if strings.Contains(inFile2, ".mxproject") {
		mxprojectFile = inFile2
	}

	if mxprojectFile != "" {
		mxprojectAll, _ := stm32cubemx.IniReader(mxprojectFile, false)

		if params.Board == "" && params.Device == "" {
			params.Board = "Test Board"
			params.Device = "Test Device"
		}

		var parms cbuild.ParamsType

		err := stm32cubemx.ReadCbuildYmlFile(inFile, outPath, &parms)
		if err != nil {
			return err
		}
		workDir := path.Dir(inFile)
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

		fPaths, err := stm32cubemx.ReadContexts(workDir+"/STM32CubeMX/STM32CubeMX.ioc", parms)
		if err != nil {
			return err
		}

		err = stm32cubemx.WriteCgenYml(outPath, mxprojectAll, fPaths, params)
		if err != nil {
			return err
		}
	} else {
		return errors.New("input file not supported")
	}

	return nil
}
