/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"path"

	"github.com/open-cmsis-pack/generator-bridge/cmd/common"
	log "github.com/sirupsen/logrus"
)

type Params_s struct {
	Board  string
	Device string
	//Packs  []struct {
	//	Pack string
	//	Path string
	//}
}

// https://zhwt.github.io/yaml-to-go/

type CbuildGenIdx_s struct {
	BuildGenIdx struct {
		GeneratedBy string `yaml:"generated-by"`
		Generators  []struct {
			ID          string `yaml:"id"`
			Device      string `yaml:"device"`
			ProjectType string `yaml:"project-type"`
			CbuildGens  []struct {
				CbuildGen     string `yaml:"cbuild-gen"`
				Project       string `yaml:"project"`
				Configuration string `yaml:"configuration"`
			} `yaml:"cbuild-gens"`
		} `yaml:"generators"`
	} `yaml:"build-gen-idx"`
}

type CbuildGen_S struct {
	BuildGen struct {
		GeneratedBy string `yaml:"generated-by"`
		Solution    string `yaml:"solution"`
		Project     string `yaml:"project"`
		Context     string `yaml:"context"`
		Compiler    string `yaml:"compiler"`
		Device      string `yaml:"device"`
		Packs       []struct {
			Pack string `yaml:"pack"`
			Path string `yaml:"path"`
		} `yaml:"packs"`
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
			FromPack   string `yaml:"from-pack"`
			SelectedBy string `yaml:"selected-by"`
		} `yaml:"components"`
		Linker struct {
			Script  string `yaml:"script"`
			Regions string `yaml:"regions"`
		} `yaml:"linker"`
		ConstructedFiles []struct {
			File     string `yaml:"file"`
			Category string `yaml:"category"`
		} `yaml:"constructed-files"`
		Licenses []struct {
			License string `yaml:"license"`
			Packs   []struct {
				Pack string `yaml:"pack"`
			} `yaml:"packs"`
			Components []struct {
				Component string `yaml:"component"`
			} `yaml:"components"`
		} `yaml:"licenses"`
	} `yaml:"build-gen"`
}

func Read(name string, params *Params_s) error {
	var cbuildGenIdx CbuildGenIdx_s

	common.ReadYml(name, &cbuildGenIdx)

	for idGen := range cbuildGenIdx.BuildGenIdx.Generators {
		generator := cbuildGenIdx.BuildGenIdx.Generators[idGen]
		genId := generator.ID
		genDevice := generator.Device
		genType := generator.ProjectType

		log.Infof("Found Generator: id: %v, device: %v, type: %v", genId, genDevice, genType)

		params.Board = ""
		params.Device = genDevice
		//split := strings.SplitAfter(cbuildGen.BuildGen.Board, "::")
		//if len(split) == 2 {
		//	params.Board = split[1]
		//} else {
		//	params.Board = cbuildGen.Build.Board
		//}

		for idSub := range generator.CbuildGens {
			cbuildGen := generator.CbuildGens[idSub]
			fileName := cbuildGen.CbuildGen
			subPath := path.Join(path.Dir(name), fileName)
			ReadCore(subPath, params)
		}
	}

	return nil
}

func ReadCore(name string, params *Params_s) error {
	var cbuildGen CbuildGen_S

	common.ReadYml(name, &cbuildGen)

	//for p := range cbuildGen.BuildGen.Packs {
	//	params.Packs = append(params.Packs, struct {
	//		Pack string
	//		Path string
	//	}(cbuildGen.BuildGen.Packs[p]))
	//}

	return nil
}
