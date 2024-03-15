/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package cbuild

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/common"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

type PackType struct {
	Pack string
	Path string
}

type SubsystemIdxType struct {
	Project              string
	CbuildGen            string
	Configuration        string
	ForProjectPart       string
	ProjectType          string
	SecureContextName    string
	NonSecureContextName string
}

type SubsystemType struct {
	SubsystemIdx SubsystemIdxType
	Board        string
	Device       string
	Project      string
	Compiler     string
	TrustZone    string
	CoreName     string
	// Packs        []PackType
}

type ParamsType struct {
	Board       string
	Device      string
	OutPath     string
	ProjectType string
	Subsystem   []SubsystemType
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

func Read(name, outPath string, params *ParamsType) error {
	return ReadCbuildgenIdx(name, outPath, params)
}

func ReadCbuildgenIdx(name, outPath string, params *ParamsType) error {
	var cbuildGenIdx CbuildGenIdxType

	err := common.ReadYml(name, &cbuildGenIdx)
	if err != nil {
		return err
	}

	for idGen, cbuildGenIdx := range cbuildGenIdx.BuildGenIdx.Generators {
		cbuildGenIdxID := cbuildGenIdx.ID
		cbuildGenIdxBoard := cbuildGenIdx.Board
		cbuildGenIdxDevice := cbuildGenIdx.Device
		cbuildGenIdxType := cbuildGenIdx.ProjectType
		cbuildGenIdxOutputPath := cbuildGenIdx.Output

		log.Infof("Found CBuildGenIdx: #%v ID: %v, board: %v, device: %v, type: %v", idGen, cbuildGenIdxID, cbuildGenIdxBoard, cbuildGenIdxDevice, cbuildGenIdxType)
		log.Infof("CBuildGenIdx Output path: %v", cbuildGenIdxOutputPath)

		params.Device = cbuildGenIdxDevice
		params.OutPath = cbuildGenIdxOutputPath

		split := strings.SplitAfter(cbuildGenIdx.Board, "::")
		if len(split) == 2 {
			params.Board = split[1]
		} else {
			params.Board = cbuildGenIdx.Board
		}

		var secureContextName string
		var nonsecureContextName string

		for _, cbuildGen := range cbuildGenIdx.CbuildGens {
			fileName := cbuildGen.CbuildGen
			var subPath string
			if filepath.IsAbs(fileName) {
				subPath = fileName
			} else {
				subPath = path.Join(path.Dir(name), fileName)
			}

			var subsystem SubsystemType
			subsystem.SubsystemIdx.Project = cbuildGen.Project
			subsystem.SubsystemIdx.Configuration = cbuildGen.Configuration
			subsystem.SubsystemIdx.CbuildGen = cbuildGen.CbuildGen
			subsystem.SubsystemIdx.ProjectType = cbuildGenIdx.ProjectType
			subsystem.SubsystemIdx.ForProjectPart = cbuildGen.ForProjectPart

			err := ReadCbuildgen(subPath, &subsystem) // use copy, do not override for next instance
			if err != nil {
				return err
			}

			params.Subsystem = append(params.Subsystem, subsystem)

			// store Reference project for TZ
			if cbuildGenIdx.ProjectType == "trustzone" {
				if cbuildGen.ForProjectPart == "secure" {
					secureContextName = cbuildGen.Project
				} else if cbuildGen.ForProjectPart == "non-secure" {
					nonsecureContextName = cbuildGen.Project
				}
			}
		}

		// store Reference project for TZ-NS
		for idSub := range params.Subsystem {
			subsystem := &params.Subsystem[idSub]
			subsystem.SubsystemIdx.SecureContextName = secureContextName
			subsystem.SubsystemIdx.NonSecureContextName = nonsecureContextName
		}
	}

	return nil
}

func ReadCbuildgen(name string, subsystem *SubsystemType) error {
	var cbuildGen CbuildGenType

	if !utils.FileExists(name) {
		text := "File not found: "
		text += name
		return errors.New(text)
	}

	err := common.ReadYml(name, &cbuildGen)
	if err != nil {
		return err
	}

	split := strings.SplitAfter(cbuildGen.BuildGen.Board, "::")
	if len(split) == 2 {
		subsystem.Board = split[1]
	} else {
		subsystem.Board = cbuildGen.BuildGen.Board
	}
	subsystem.Device = cbuildGen.BuildGen.Device
	subsystem.Compiler = cbuildGen.BuildGen.Compiler
	subsystem.Project = cbuildGen.BuildGen.Project
	subsystem.CoreName = cbuildGen.BuildGen.Processor.Core
	subsystem.TrustZone = cbuildGen.BuildGen.Processor.Trustzone

	log.Infof("Found CBuildGen: board: %v, device: %v, core: %v, TZ: %v, compiler: %v, project: %v",
		subsystem.Board, subsystem.Device, subsystem.CoreName, subsystem.TrustZone, subsystem.Compiler, subsystem.Project)

	// for id := range cbuildGen.BuildGen.Packs {
	// 	genPack := cbuildGen.BuildGen.Packs[id]
	// 	var pack PackType
	// 	pack.Pack = genPack.Pack
	// 	pack.Path = genPack.Path
	// 	log.Infof("Found Pack: #%v Pack: %v, Path: %v", id, pack.Pack, pack.Path)
	// 	subsystem.Packs = append(subsystem.Packs, pack)
	// }

	return nil
}
