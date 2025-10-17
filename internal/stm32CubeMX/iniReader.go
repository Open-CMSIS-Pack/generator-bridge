/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/ini.v1"
)

type MxprojectAllType struct {
	Mxproject []MxprojectType
}

type MxprojectType struct {
	Context          string
	PreviousLibFiles struct {
		LibFiles []string
	}
	PreviousUsedFiles struct {
		SourceFiles []string
		HeaderPath  []string
		CDefines    []string
	}
	PreviousGenFiles struct {
		AdvancedFolderStructure string
		HeaderFilesList         []string
		HeaderPathList          []string
		HeaderFiles             string
		SourceFilesList         []string
		SourcePathList          []string
		SourceFiles             string
	}
	ThirdPartyIpFiles struct {
		IncludeFiles   []string
		SourceAsmFiles []string
		SourceFiles    []string
	}
}

type IniSectionsType struct {
	Context string
	Section string
}

func PrintKeyValStr(key, val string) {
	fmt.Printf("\n%v : %v", key, val)
}

func PrintKeyValStrs(key string, vals []string) {
	fmt.Printf("\n%v", key)
	for i := range vals {
		fmt.Printf("\n%d: %v", i, vals[i])
	}
}

func PrintKeyValInt(key string, val int) {
	fmt.Printf("\n%v : %v", key, val)
}

func PrintItemCsv(section *ini.Section, key string) {
	valStr := section.Key(key).String()
	commentStr := section.Key(key).Comment
	commentStrs := strings.Split(commentStr, ";")
	PrintKeyValStr(key, valStr)
	PrintKeyValStrs(key, commentStrs)
}

func PrintItem(section *ini.Section, key string) {
	valStr := section.Key(key).String()
	PrintKeyValStr(key, valStr)
}

func PrintItemIterator(section *ini.Section, key, iterator string) {
	valStr := section.Key(key).String()
	PrintKeyValStr(key, valStr)
	maxCnt, _ := strconv.Atoi(valStr)
	for cnt := 0; cnt < maxCnt; cnt++ {
		keyN := iterator + strconv.Itoa(cnt)
		valStr = section.Key(keyN).String()
		PrintKeyValStr(keyN, valStr)
	}
}

func StoreData(data *string, value string) {
	value = filepath.ToSlash(value)

	if value != "" {
		*data = value
	}
}

func StoreDataArray(data *[]string, values ...string) {
	for _, value := range values {
		value = filepath.ToSlash(value)
		if value != "" {
			if !slices.Contains(*data, value) {
				*data = append(*data, value)
			}
		}
	}
}

func StoreItem(data *string, section *ini.Section, key string) {
	StoreData(data, section.Key(key).String())
}

func StoreItemCsv(data *[]string, section *ini.Section, key string) {
	valStr := section.Key(key).String()
	commentStr := section.Key(key).Comment
	commentStrs := strings.Split(commentStr, ";")
	StoreDataArray(data, valStr)
	StoreDataArray(data, commentStrs...)
}

func StoreItemIterator(data *[]string, section *ini.Section, key, iterator string) {
	valStr := section.Key(key).String()
	maxCnt, _ := strconv.Atoi(valStr)
	for cnt := 0; cnt < maxCnt; cnt++ {
		keyN := iterator + strconv.Itoa(cnt)
		valStr = section.Key(keyN).String()
		StoreDataArray(data, valStr)
	}
}

func IniReader(path string, params []BridgeParamType) (MxprojectAllType, error) {
	var mxprojectAll MxprojectAllType

	if !utils.FileExists(path) {
		text := "File not found: "
		text += path
		return mxprojectAll, errors.New(text)
	}

	inidata, err := GetIni(path)
	if err != nil || inidata == nil {
		text := "File not found or error opening file: .mxproject"
		text += path
		return mxprojectAll, errors.New(text)
	}

	var iniSections []IniSectionsType
	err = GetSections(inidata, &iniSections)
	if err != nil {
		return mxprojectAll, err
	}

	for _, param := range params {
		context := param.CubeContext

		mxproject, _ := GetData(inidata, context, param.Compiler)
		mxproject.Context = context
		mxprojectAll.Mxproject = append(mxprojectAll.Mxproject, mxproject)
	}

	return mxprojectAll, nil
}

func GetIni(path string) (*ini.File, error) {
	log.Debugf("\nReading CubeMX config file: %v", path)

	inidata, err := ini.Load(path)
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return inidata, nil
	}

	return inidata, nil
}

func GetSections(inidata *ini.File, iniSections *[]IniSectionsType) error {
	sectionsData := inidata.SectionStrings()
	for _, section := range sectionsData {
		var iniSection IniSectionsType
		sectionString := strings.Split(section, ":")
		if len(sectionString) > 1 {
			iniSection.Context = sectionString[0]
			iniSection.Section = sectionString[1]
		} else {
			iniSection.Context = ""
			iniSection.Section = section
		}

		*iniSections = append(*iniSections, iniSection)
	}

	return nil
}

func GetData(inidata *ini.File, iniName string, compiler string) (MxprojectType, error) {
	var mxproject MxprojectType
	var sectionName string
	var PreviousUsedFilesID string
	var section *ini.Section

	const ThirdPartyIpID = "ThirdPartyIp"
	if iniName != "" {
		sectionName = iniName + ":" + ThirdPartyIpID
	} else {
		sectionName = ThirdPartyIpID
	}
	section = inidata.Section(sectionName)
	if section != nil {
		var ipNames []string
		ipNumber := section.Key("ThirdPartyIpNumber").String()
		ipCnt, _ := strconv.Atoi(ipNumber)
		for cnt := 0; cnt < ipCnt; cnt++ {
			ipName := section.Key("ThirdPartyIpName#" + strconv.Itoa(cnt)).String()
			if ipName != "" {
				ipNames = append(ipNames, ipName)
			}
		}
		for _, ipName := range ipNames {
			ipName = "ThirdPartyIp#" + ipName
			if iniName != "" {
				sectionName = iniName + ":" + ipName
			} else {
				sectionName = ipName
			}
			section = inidata.Section(sectionName)
			if section != nil {
				StoreItemCsv(&mxproject.ThirdPartyIpFiles.IncludeFiles, section, "include")
				StoreItemCsv(&mxproject.ThirdPartyIpFiles.SourceAsmFiles, section, "sourceAsm")
				StoreItemCsv(&mxproject.ThirdPartyIpFiles.SourceFiles, section, "source")
			}
		}
	}

	PreviousUsedFilesID, err := GetPreviousUsedFilesID(compiler)
	if err != nil {
		return mxproject, err
	}

	if iniName != "" {
		sectionName = iniName + ":" + PreviousUsedFilesID
	} else {
		sectionName = PreviousUsedFilesID
	}
	section = inidata.Section(sectionName)
	if section != nil {
		StoreItemCsv(&mxproject.PreviousUsedFiles.SourceFiles, section, "SourceFiles")
		StoreItemCsv(&mxproject.PreviousUsedFiles.HeaderPath, section, "HeaderPath")
		StoreItemCsv(&mxproject.PreviousUsedFiles.CDefines, section, "CDefines")
		PrintItemCsv(section, "SourceFiles")
		PrintItemCsv(section, "HeaderPath")
		PrintItemCsv(section, "CDefines")
	}

	const PreviousLibFilesID = "PreviousLibFiles"
	if iniName != "" {
		sectionName = iniName + ":" + PreviousLibFilesID
	} else {
		sectionName = PreviousLibFilesID
	}
	section = inidata.Section(sectionName)
	if section != nil {
		StoreItemCsv(&mxproject.PreviousLibFiles.LibFiles, section, "LibFiles")
		PrintItemCsv(section, "LibFiles")
	}

	const PreviousGenFilesID = "PreviousGenFiles"
	if iniName != "" {
		sectionName = iniName + ":" + PreviousGenFilesID
	} else {
		sectionName = PreviousGenFilesID
	}
	section = inidata.Section(sectionName)
	if section != nil {
		StoreItem(&mxproject.PreviousGenFiles.AdvancedFolderStructure, section, "AdvancedFolderStructure")
		StoreItemIterator(&mxproject.PreviousGenFiles.HeaderFilesList, section, "HeaderFileListSize", "HeaderFiles#")
		StoreItemIterator(&mxproject.PreviousGenFiles.HeaderPathList, section, "HeaderFolderListSize", "HeaderPath#")
		StoreItem(&mxproject.PreviousGenFiles.HeaderFiles, section, "HeaderFiles")
		StoreItemIterator(&mxproject.PreviousGenFiles.SourceFilesList, section, "SourceFileListSize", "SourceFiles#")
		StoreItemIterator(&mxproject.PreviousGenFiles.HeaderFilesList, section, "HeaderFileListSize", "HeaderFiles#")
		StoreItemIterator(&mxproject.PreviousGenFiles.SourcePathList, section, "SourceFolderListSize", "SourcePath#")
		StoreItem(&mxproject.PreviousGenFiles.SourceFiles, section, "SourceFiles")

		PrintItem(section, "AdvancedFolderStructure")
		PrintItemIterator(section, "HeaderFileListSize", "HeaderFiles#")
		PrintItemIterator(section, "HeaderFolderListSize", "HeaderPath#")
		PrintItem(section, "HeaderFiles")
		PrintItemIterator(section, "SourceFileListSize", "SourceFiles#")
		PrintItemIterator(section, "HeaderFileListSize", "HeaderFiles#")
		PrintItemIterator(section, "SourceFolderListSize", "SourcePath#")
		PrintItem(section, "SourceFiles")
	}

	return mxproject, nil
}

func GetPreviousUsedFilesID(compiler string) (string, error) {
	var sectionMapping = map[string]string{
		"AC6":   "PreviousUsedKeilFiles",
		"GCC":   "PreviousUsedCubeIDEFiles",
		"IAR":   "PreviousUsedIarFiles",
		"CLANG": "PreviousUsedCubeIDEFiles",
	}

	PreviousUsedFilesID, ok := sectionMapping[compiler]
	if !ok {
		return "", errors.New("unknown compiler")
	}
	return PreviousUsedFilesID, nil
}
