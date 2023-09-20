/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/cmd/common"
)

type Params_s struct {
	Board  string
	Device string
	Packs  []struct {
		Pack string
		Path string
	}
}

// https://zhwt.github.io/yaml-to-go/
type Conf_s struct {
	Build struct {
		GeneratedBy string `yaml:"generated-by"`
		Solution    string `yaml:"solution"`
		Project     string `yaml:"project"`
		Context     string `yaml:"context"`
		Compiler    string `yaml:"compiler"`
		Board       string `yaml:"board"`
		Device      string `yaml:"device"`
		Processor   struct {
			Fpu       string `yaml:"fpu"`
			Endian    string `yaml:"endian"`
			Trustzone string `yaml:"trustzone"`
		} `yaml:"processor"`
		Packs []struct {
			Pack string `yaml:"pack"`
			Path string `yaml:"path"`
		} `yaml:"packs"`
		Optimize string `yaml:"optimize"`
		Debug    string `yaml:"debug"`
		Misc     struct {
			ASM  []string `yaml:"ASM"`
			C    []string `yaml:"C"`
			CPP  []string `yaml:"CPP"`
			Link []string `yaml:"Link"`
		} `yaml:"misc"`
		Define     []string `yaml:"define"`
		AddPath    []string `yaml:"add-path"`
		OutputDirs struct {
			Intdir string `yaml:"intdir"`
			Outdir string `yaml:"outdir"`
			Rtedir string `yaml:"rtedir"`
		} `yaml:"output-dirs"`
		Output []struct {
			Type string `yaml:"type"`
			File string `yaml:"file"`
		} `yaml:"output"`
		Components []struct {
			Component  string `yaml:"component"`
			Condition  string `yaml:"condition"`
			FromPack   string `yaml:"from-pack"`
			SelectedBy string `yaml:"selected-by"`
			Files      []struct {
				File     string `yaml:"file"`
				Category string `yaml:"category"`
			} `yaml:"files"`
			Generator struct {
				ID string `yaml:"id"`
			} `yaml:"generator"`
		} `yaml:"components"`
		Generators []struct {
			Generator string `yaml:"generator"`
			Path      string `yaml:"path"`
			Gpdsc     string `yaml:"gpdsc"`
			Command   struct {
				Win struct {
					File      string   `yaml:"file"`
					Arguments []string `yaml:"arguments"`
				} `yaml:"win"`
				Linux struct {
					File      string   `yaml:"file"`
					Arguments []string `yaml:"arguments"`
				} `yaml:"linux"`
				Mac struct {
					File      string   `yaml:"file"`
					Arguments []string `yaml:"arguments"`
				} `yaml:"mac"`
				Other struct {
					File      string   `yaml:"file"`
					Arguments []string `yaml:"arguments"`
				} `yaml:"other"`
			} `yaml:"command"`
		} `yaml:"generators"`
		Linker struct {
			Script  string `yaml:"script"`
			Regions string `yaml:"regions"`
		} `yaml:"linker"`
		ConstructedFiles []struct {
			File     string `yaml:"file"`
			Category string `yaml:"category"`
		} `yaml:"constructed-files"`
	} `yaml:"build"`
}

func Read(path string, params *Params_s) error {
	var conf Conf_s

	common.ReadYml(path, &conf)
	split := strings.SplitAfter(conf.Build.Board, "::")

	if len(split) == 2 {
		params.Board = split[1]
	} else {
		params.Board = conf.Build.Board
	}
	params.Device = conf.Build.Device

	for p := range conf.Build.Packs {
		params.Packs = append(params.Packs, struct {
			Pack string
			Path string
		}(conf.Build.Packs[p]))
	}

	return nil
}
