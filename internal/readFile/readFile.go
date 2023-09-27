/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package readFile

import (
	"errors"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	log "github.com/sirupsen/logrus"
)

func Process(inFile, outPath string) error {
	log.Infof("Reading file: %v", inFile)

	if strings.Contains(inFile, "cbuild-gen-idx.yml") {
		var params cbuild.ParamsType
		cbuild.Read(inFile, outPath, &params)
	} else if strings.Contains(inFile, "cbuild-gen.yml") || strings.Contains(inFile, "cbuild.yml") {
		var params cbuild.ParamsType
		params.OutPath = outPath
		cbuild.ReadCbuildgen(inFile, &params)
	} else if strings.Contains(inFile, ".mxproject") {
		mxproject, _ := stm32CubeMX.IniReader(inFile, false)

		var inParms cbuild.ParamsType
		inParms.Board = "Test Board"
		inParms.Device = "Test Device"

		if outPath != "" {
			stm32CubeMX.WriteCgenYml(outPath, mxproject, inParms)
		}
	} else {
		return errors.New("input file not supported")
	}

	return nil
}
