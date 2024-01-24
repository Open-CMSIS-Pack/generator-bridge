/*
 * Copyright (c) 2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package generator

import (
	"errors"

	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
)

type ParamsType struct {
	ID          string
	DownloadURL string
}

type GeneratorType struct {
	Generator []struct {
		ID          string `yaml:"id"`
		Description string `yaml:"description"`
		DownloadURL string `yaml:"download-url"`
		Run         string `yaml:"run"`
		Path        string `yaml:"path"`
	} `yaml:"generator"`
}

func Read(name string, params *ParamsType) error {
	var gen GeneratorType

	if !utils.FileExists(name) {
		text := "File not found: "
		text += name
		return errors.New(text)
	}

	err := common.ReadYml(name, &gen)
	if err != nil {
		return err
	}
	for _, genx := range gen.Generator {
		if genx.ID == "CubeMX" {
			params.ID = genx.ID
			params.DownloadURL = genx.DownloadURL
			break
		}
	}
	if params.ID != "CubeMX" {
		return errors.New("generator CubeMX missing in global.generator.yml")
	}
	return nil
}
