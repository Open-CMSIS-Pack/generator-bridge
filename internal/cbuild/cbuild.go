/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"path"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	log "github.com/sirupsen/logrus"
)

type PackType struct {
	Pack string
	Path string
}

type CoreType struct {
	Board    string
	Device   string
	Project  string
	Compiler string
	Packs    []PackType
}

type ParamsType struct {
	Board   string
	Device  string
	OutPath string
	Core    []CoreType
}

// https://zhwt.github.io/yaml-to-go/

// IDX input file
type CbuildGenIdxType struct {
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
				Output        string `yaml:"output"`
			} `yaml:"cbuild-gens"`
		} `yaml:"generators"`
	} `yaml:"build-gen-idx"`
}

// Sub input file
type CbuildGenType struct {
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
		Groups []struct {
			Group string `yaml:"group"`
			Files []struct {
				File     string `yaml:"file"`
				Category string `yaml:"category"`
			} `yaml:"files"`
		} `yaml:"groups"`
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

// bridge generator output file
type CgenType struct {
	Layer LayerType `yaml:"layer"`
}
type CgenPacksType struct {
	Pack string `yaml:"pack"`
}
type CgenFilesType struct {
	File string `yaml:"file"`
}
type CgenGroupsType struct {
	Group string          `yaml:"group"`
	Files []CgenFilesType `yaml:"files"`
}
type LayerType struct {
	ForDevice string           `yaml:"for-device,omitempty"`
	ForBoard  string           `yaml:"for-board,omitempty"`
	Packs     []CgenPacksType  `yaml:"packs,omitempty"` // do not set if no new packs
	Define    []string         `yaml:"define,omitempty"`
	AddPath   []string         `yaml:"add-path,omitempty"`
	Groups    []CgenGroupsType `yaml:"groups"`
}

func Read(name, outPath string, params *ParamsType) error {
	return ReadCbuildgenIdx(name, outPath, params)
}

func ReadCbuildgenIdx(name, outPath string, params *ParamsType) error {
	var cbuildGenIdx CbuildGenIdxType

	err := common.ReadYml(name, &cbuildGenIdx)
	if err != nil {
		return err
	}

	for idGen := range cbuildGenIdx.BuildGenIdx.Generators {
		cbuildGenIdx := cbuildGenIdx.BuildGenIdx.Generators[idGen]
		cbuildGenIdxID := cbuildGenIdx.ID
		cbuildGenIdxBoard := cbuildGenIdx.Board
		cbuildGenIdxDevice := cbuildGenIdx.Device
		cbuildGenIdxType := cbuildGenIdx.ProjectType

		log.Infof("Found CBuildGenIdx: #%v Id: %v, board: %v, device: %v, type: %v", idGen, cbuildGenIdxID, cbuildGenIdxBoard, cbuildGenIdxDevice, cbuildGenIdxType)

		params.Device = cbuildGenIdxDevice
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

			var outputPath string
			if cbuildGen.Output != "" {
				params.OutPath = cbuildGen.Output
			} else {
				params.OutPath = outputPath
			}

			err := ReadCbuildgen(subPath, params) // use copy, do not override for next instance
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadCbuildgen(name string, params *ParamsType) error {
	var cbuildGen CbuildGenType

	err := common.ReadYml(name, &cbuildGen)
	if err != nil {
		return err
	}

	var core CoreType

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
		var pack PackType
		pack.Pack = genPack.Pack
		pack.Path = genPack.Path
		log.Infof("Found Pack: #%v Pack: %v, Path: %v", id, pack.Pack, pack.Path)
		core.Packs = append(core.Packs, pack)
	}

	params.Core = append(params.Core, core)

	return nil
}
