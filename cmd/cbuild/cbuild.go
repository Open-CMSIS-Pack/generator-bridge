/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"path"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/cmd/common"
	log "github.com/sirupsen/logrus"
)

type Pack_s struct {
	Pack string
	Path string
}

type Core_s struct {
	Board    string
	Device   string
	Project  string
	Compiler string
	Packs    []Pack_s
}

type Params_s struct {
	Board  string
	Device string
	Core   []Core_s
}

// https://zhwt.github.io/yaml-to-go/

type CbuildGenIdx_s struct {
	BuildGenIdx struct {
		GeneratedBy string `yaml:"generated-by"`
		Generators  []struct {
			ID          string `yaml:"id"`
			Device      string `yaml:"device"`
			Board       string `yaml:"board"`
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
		Board       string `yaml:"board"`
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

type Cgen_s struct {
	Generator CgenGenerator_s `yaml:"Generator"`
}
type CgenPacks_s struct {
	Pack string `yaml:"pack"`
}
type CgenFiles_s struct {
	File string `yaml:"file"`
}
type CgenGroups_s struct {
	Group string        `yaml:"group"`
	Files []CgenFiles_s `yaml:"files"`
}
type CgenGenerator_s struct {
	GeneratedBy string         `yaml:"generated-by"`
	ForDevice   string         `yaml:"for-device"`
	ForBoard    string         `yaml:"for-board"`
	Packs       []CgenPacks_s  `yaml:"packs"`
	Define      []string       `yaml:"define"`
	AddPath     []string       `yaml:"add-path"`
	Groups      []CgenGroups_s `yaml:"groups"`
}

func Read(name string, params *Params_s) error {
	return ReadCbuildgenIdx(name, params)
}

func ReadCbuildgenIdx(name string, params *Params_s) error {
	var cbuildGenIdx CbuildGenIdx_s

	common.ReadYml(name, &cbuildGenIdx)

	for idGen := range cbuildGenIdx.BuildGenIdx.Generators {
		cbuildGenIdx := cbuildGenIdx.BuildGenIdx.Generators[idGen]
		genId := cbuildGenIdx.ID
		genBoard := cbuildGenIdx.Board
		genDevice := cbuildGenIdx.Device
		genType := cbuildGenIdx.ProjectType

		log.Infof("Found CBuildGenIdx: #%v Id: %v, board: %v, device: %v, type: %v", idGen, genId, genBoard, genDevice, genType)

		params.Device = genDevice
		split := strings.SplitAfter(cbuildGenIdx.Board, "::")
		if len(split) == 2 {
			params.Board = split[1]
		} else {
			params.Board = cbuildGenIdx.Board
		}

		for idSub := range cbuildGenIdx.CbuildGens {
			cbuildGen := cbuildGenIdx.CbuildGens[idSub]
			fileName := cbuildGen.CbuildGen
			subPath := path.Join(path.Dir(name), fileName)
			ReadCbuildgen(subPath, params)
		}
	}

	return nil
}

func ReadCbuildgen(name string, params *Params_s) error {
	var cbuildGen CbuildGen_S

	common.ReadYml(name, &cbuildGen)
	var core Core_s

	split := strings.SplitAfter(cbuildGen.BuildGen.Board, "::")
	if len(split) == 2 {
		core.Board = split[1]
	} else {
		core.Board = cbuildGen.BuildGen.Board
	}
	core.Device = cbuildGen.BuildGen.Device
	core.Compiler = cbuildGen.BuildGen.Compiler
	core.Project = cbuildGen.BuildGen.Project

	log.Infof("Found CBuildGen: board: %v, device: %v, compiler: %v, project: %v", core.Board, core.Device, core.Compiler, core.Project)

	for id := range cbuildGen.BuildGen.Packs {
		genPack := cbuildGen.BuildGen.Packs[id]
		var pack Pack_s
		pack.Pack = genPack.Pack
		pack.Path = genPack.Path
		log.Infof("Found Pack: #%v Pack: %v, Path: %v", id, pack.Pack, pack.Path)
		core.Packs = append(core.Packs, pack)
	}

	params.Core = append(params.Core, core)

	return nil
}
