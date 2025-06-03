/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package readfile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	stm32cubemx "github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	log "github.com/sirupsen/logrus"
)

func Process(inFile, inFile2, outPath string) error {
	log.Debugf("Reading file: %v", inFile)
	if outPath == "" {
		outPath = filepath.Dir(inFile2)
	}

	var cbuildParams cbuild.ParamsType
	var params []stm32cubemx.BridgeParamType

	if strings.Contains(inFile, "cbuild-gen-idx.yml") {
		err := cbuild.Read(inFile, "CubeMX", &cbuildParams)
		if err != nil {
			return err
		}

		err = stm32cubemx.GetBridgeInfo(&cbuildParams, &params)
		if err != nil {
			return err
		}

		_, err = stm32cubemx.WriteProjectFile(outPath, params[0])
		if err != nil {
			return err
		}
	}

	var mxprojectFile string
	if strings.Contains(inFile, ".mxproject") {
		mxprojectFile = inFile
	}
	if strings.Contains(inFile2, ".mxproject") {
		mxprojectFile = inFile2
	}

	if mxprojectFile != "" {
		mxprojectAll, _ := stm32cubemx.IniReader(mxprojectFile, params)

		if params[0].BoardName == "" && params[0].Device == "" {
			params[0].BoardName = "Test Board"
			params[0].Device = "Test Device"
		}

		err := stm32cubemx.ReadCbuildGenIdxYmlFile(inFile, "CubeMX", &cbuildParams)
		if err != nil {
			return err
		}
		workDir := filepath.Dir(inFile)
		if params[0].Output != "" {
			if filepath.IsAbs(params[0].Output) {
				workDir = params[0].Output
			} else {
				workDir = filepath.Join(workDir, params[0].Output)
			}
		} else {
			if filepath.IsAbs(outPath) {
				workDir = outPath
			} else {
				workDir = filepath.Join(workDir, outPath)
			}
		}
		workDir = filepath.Clean(workDir)
		workDir = filepath.ToSlash(workDir)
		err = os.MkdirAll(workDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = stm32cubemx.ReadContexts(workDir+"/STM32CubeMX/STM32CubeMX.ioc", params)
		if err != nil {
			return err
		}

		err = stm32cubemx.WriteCgenYml(outPath, mxprojectAll, params)
		if err != nil {
			return err
		}
	} else {
		return errors.New("input file not supported")
	}

	return nil
}
