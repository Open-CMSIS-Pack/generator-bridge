/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/ini.v1"
)

type MxprojectAllType struct {
	Mxproject []MxprojectType
}

type MxprojectType struct {
	Pname            string
	Trustzone        string
	PreviousLibFiles struct {
		LibFiles []string
	}
	PreviousUsedKeilFiles struct {
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
}

type IniSectionCore struct {
	pname     string
	trustzone string
	iniName   string
}

type IniSectionsType struct {
	cores    []IniSectionCore
	sections []string
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
	for id := range values {
		value := values[id]
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

func FindInList(name string, list *[]string) bool {
	if name == "" {
		return false
	}

	for id := range *list {
		item := (*list)[id]
		if item == name {
			return true
		}
	}
	return false
}

func AppendToList(name string, list *[]string) {
	if name == "" {
		return
	}

	if FindInList(name, list) {
		return
	}

	*list = append(*list, name)
}

func FindInCores(name string, list *[]IniSectionCore) bool {
	if name == "" {
		return false
	}

	for id := range *list {
		item := (*list)[id]
		if item.iniName == name {
			return true
		}
	}
	return false
}

func AppendToCores(iniSectionCore IniSectionCore, list *[]IniSectionCore) {
	name := iniSectionCore.iniName
	if name == "" {
		return
	}

	if FindInCores(name, list) {
		return
	}

	*list = append(*list, iniSectionCore)
}

func IniReader(path string, trustzone bool) (MxprojectAllType, error) {
	var mxprojectAll MxprojectAllType
	inidata, err := GetIni(path)
	if err != nil {
		return mxprojectAll, err
	}

	var iniSections IniSectionsType
	err = GetSections(inidata, &iniSections)
	if err != nil {
		return mxprojectAll, err
	}

	if len(iniSections.cores) == 0 { // single core
		sectionName := ""
		mxproject, _ := GetData(inidata, sectionName)
		mxprojectAll.Mxproject = append(mxprojectAll.Mxproject, mxproject)
	} else { // multi core / trustzone
		for coreId := range iniSections.cores {
			core := iniSections.cores[coreId]
			iniName := core.iniName
			pname := core.pname
			trustzone := core.trustzone
			mxproject, _ := GetData(inidata, iniName)
			mxproject.Pname = pname
			mxproject.Trustzone = trustzone
			mxprojectAll.Mxproject = append(mxprojectAll.Mxproject, mxproject)
		}
	}

	return mxprojectAll, nil
}

func GetIni(path string) (*ini.File, error) {
	log.Infof("\nReading CubeMX config file: %v", path)

	inidata, err := ini.Load(path)
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return inidata, nil
	}

	return inidata, nil
}

func GetSections(inidata *ini.File, iniSections *IniSectionsType) error {
	sectionsData := inidata.SectionStrings()
	for sectionId := range sectionsData {
		section := sectionsData[sectionId]
		var iniName string
		var sectionName string
		sectionString := strings.Split(section, ":")
		if len(sectionString) > 1 {
			iniName = sectionString[0]
			sectionName = sectionString[1]
		} else {
			iniName = ""
			sectionName = section
		}

		var pname string
		re := regexp.MustCompile("[0-9]+")
		coreNameNumbers := re.FindAllString(iniName, -1)
		if len(coreNameNumbers) == 1 {
			pname = "Cortex-M" + coreNameNumbers[0]
		}

		var trustzone string
		iniLen := len(iniName)
		if iniLen > 0 {
			if strings.LastIndex(iniName, "S") == iniLen-1 {
				if strings.LastIndex(iniName, "NS") == iniLen-2 {
					trustzone = "non-secure"
				} else {
					trustzone = "secure"
				}
			}
		}

		var iniSectionCore IniSectionCore
		iniSectionCore.iniName = iniName
		iniSectionCore.pname = pname
		iniSectionCore.trustzone = trustzone
		AppendToCores(iniSectionCore, &iniSections.cores)
		AppendToList(sectionName, &iniSections.sections)
	}

	return nil
}

func GetData(inidata *ini.File, iniName string) (MxprojectType, error) {
	var mxproject MxprojectType

	var sectionName string
	const PreviousUsedKeilFilesId = "PreviousUsedKeilFiles"
	if iniName != "" {
		sectionName = iniName + ":" + PreviousUsedKeilFilesId
	} else {
		sectionName = PreviousUsedKeilFilesId
	}
	section := inidata.Section(sectionName)
	if section != nil {
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.SourceFiles, section, "SourceFiles")
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.HeaderPath, section, "HeaderPath")
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.CDefines, section, "CDefines")
		PrintItemCsv(section, "SourceFiles")
		PrintItemCsv(section, "HeaderPath")
		PrintItemCsv(section, "CDefines")
	}

	const PreviousLibFilesId = "PreviousLibFiles"
	if iniName != "" {
		sectionName = iniName + ":" + PreviousLibFilesId
	} else {
		sectionName = PreviousLibFilesId
	}
	section = inidata.Section(sectionName)
	if section != nil {
		StoreItemCsv(&mxproject.PreviousLibFiles.LibFiles, section, "LibFiles")
		PrintItemCsv(section, "LibFiles")
	}

	const PreviousGenFilesId = "PreviousGenFiles"
	if iniName != "" {
		sectionName = iniName + ":" + PreviousGenFilesId
	} else {
		sectionName = PreviousGenFilesId
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
