/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package readFile

import (
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/cmd/cbuild"
	"github.com/open-cmsis-pack/generator-bridge/cmd/stm32CubeMX"
	log "github.com/sirupsen/logrus"
)

func Process(inFile, outPath string) error {
	log.Infof("Reading file: %v", inFile)

	if strings.Contains(inFile, "cbuild-gen-idx.yml") {
		var params cbuild.Params_s
		cbuild.Read(inFile, &params)
	}
	if strings.Contains(inFile, "cbuild-gen.yml") {
		var params cbuild.Params_s
		cbuild.ReadCbuildgen(inFile, &params)
	}
	if strings.Contains(inFile, ".mxproject") {
		mxproject, _ := stm32CubeMX.IniReader(inFile, false)
		if outPath != "" {
			stm32CubeMX.WriteCgenYml(outPath, mxproject)
		}
	}

	return nil
}
