/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package readfile

import (
	"errors"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	stm32cubemx "github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	log "github.com/sirupsen/logrus"
)

func Process(inFile, outPath string) error {
	log.Infof("Reading file: %v", inFile)

	if strings.Contains(inFile, "cbuild-gen-idx.yml") {
		var params cbuild.ParamsType
		err := cbuild.Read(inFile, outPath, &params)
		if err != nil {
			return err
		}
	} else if strings.Contains(inFile, "cbuild-gen.yml") || strings.Contains(inFile, "cbuild.yml") {
		var params cbuild.ParamsType
		params.OutPath = outPath
		err := cbuild.ReadCbuildgen(inFile, &params)
		if err != nil {
			return err
		}
	} else if strings.Contains(inFile, ".mxproject") {
		mxproject, _ := stm32cubemx.IniReader(inFile, false)

		var inParms cbuild.ParamsType
		inParms.Board = "Test Board"
		inParms.Device = "Test Device"

		if outPath != "" {
			err := stm32cubemx.WriteCgenYml(outPath, mxproject, inParms)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("input file not supported")
	}

	return nil
}
