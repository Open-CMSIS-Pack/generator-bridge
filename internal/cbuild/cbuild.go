/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"errors"

	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
)

type CbuildGensType struct {
	CbuildGen      CbuildGenType
	Project        string
	Configuration  string
	ForProjectPart string
	Output         string
	Name           string
	Map            string
}

type ParamsType struct {
	GeneratedBy string
	ID          string
	Output      string
	Device      string
	Board       string
	ProjectType string
	CbuildGens  []CbuildGensType
}

// https://zhwt.github.io/yaml-to-go/

// IDX input file
type CbuildGenIdxType struct {
	BuildGenIdx struct {
		GeneratedBy string `yaml:"generated-by"`
		Generators  []struct {
			ID          string `yaml:"id"`
			Output      string `yaml:"output"`
			Device      string `yaml:"device"`
			Board       string `yaml:"board"`
			ProjectType string `yaml:"project-type"`
			CbuildGens  []struct {
				CbuildGen      string `yaml:"cbuild-gen"`
				Project        string `yaml:"project"`
				Configuration  string `yaml:"configuration"`
				ForProjectPart string `yaml:"for-project-part"`
				Output         string `yaml:"output"`
				Name           string `yaml:"name"`
				Map            string `yaml:"map"`
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
		Processor   struct {
			Fpu       string `yaml:"fpu"`
			Endian    string `yaml:"endian"`
			Trustzone string `yaml:"trustzone"`
			Core      string `yaml:"core"` // Dcore aus pdsc
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
			Condition  string `yaml:"condition,omitempty"`
			FromPack   string `yaml:"from-pack"`
			SelectedBy string `yaml:"selected-by"`
			Files      []struct {
				File     string `yaml:"file"`
				Category string `yaml:"category"`
				Attr     string `yaml:"attr"`
				Version  string `yaml:"version"`
			} `yaml:"files,omitempty"`
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
			License          string `yaml:"license"`
			LicenseAgreement string `yaml:"license-agreement,omitempty"`
			Packs            []struct {
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
	GeneratorImport GeneratorImportType `yaml:"generator-import,omitempty"`
}
type CgenPacksType struct {
	Pack string `yaml:"pack,omitempty"`
}
type CgenFilesType struct {
	File string `yaml:"file,omitempty"`
}
type CgenGroupsType struct {
	Group string          `yaml:"group,omitempty"`
	Files []CgenFilesType `yaml:"files,omitempty"`
}
type GeneratorImportType struct {
	ForDevice string           `yaml:"for-device,omitempty"`
	ForBoard  string           `yaml:"for-board,omitempty"`
	Packs     []CgenPacksType  `yaml:"packs,omitempty"` // do not set if no new packs
	Define    []string         `yaml:"define,omitempty"`
	AddPath   []string         `yaml:"add-path,omitempty"`
	Groups    []CgenGroupsType `yaml:"groups,omitempty"`
}

func Read(name, generatorID string, params *ParamsType) error {
	return ReadCbuildgenIdx(name, generatorID, params)
}

func ReadCbuildgenIdx(name, generatorID string, params *ParamsType) error {
	var cbuildGenIdx CbuildGenIdxType

	err := common.ReadYml(name, &cbuildGenIdx)
	if err != nil {
		return err
	}

	for _, cgen := range cbuildGenIdx.BuildGenIdx.Generators {
		if cgen.ID == generatorID {
			params.GeneratedBy = cbuildGenIdx.BuildGenIdx.GeneratedBy
			params.ID = cgen.ID
			params.Output = cgen.Output
			params.Device = cgen.Device
			params.Board = cgen.Board
			params.ProjectType = cgen.ProjectType

			for _, cbuildGen := range cgen.CbuildGens {
				var tmpCbuildGen CbuildGensType
				err := ReadCbuildgen(cbuildGen.CbuildGen, &tmpCbuildGen.CbuildGen)
				if err != nil {
					return err
				}
				tmpCbuildGen.Project = cbuildGen.Project
				tmpCbuildGen.Configuration = cbuildGen.Configuration
				tmpCbuildGen.ForProjectPart = cbuildGen.ForProjectPart
				tmpCbuildGen.Output = cbuildGen.Output
				tmpCbuildGen.Name = cbuildGen.Name
				tmpCbuildGen.Map = cbuildGen.Map

				params.CbuildGens = append(params.CbuildGens, tmpCbuildGen)
			}
		}
	}

	return nil
}

func ReadCbuildgen(name string, cbuildGen *CbuildGenType) error {

	if !utils.FileExists(name) {
		text := "File not found: "
		text += name
		return errors.New(text)
	}

	err := common.ReadYml(name, &cbuildGen)
	if err != nil {
		return err
	}
	return nil
}
