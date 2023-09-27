/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/ini.v1"
)

type MxprojectType struct {
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

func IniReader(path string, trustzone bool) (MxprojectType, error) {
	log.Infof("\nReading CubeMX config file: %v", path)

	var mxproject MxprojectType
	inidata, err := ini.Load(path)
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return mxproject, nil
	}

	section := inidata.Section("PreviousLibFiles")
	if section != nil {
		PrintItemCsv(section, "LibFiles")
		StoreItemCsv(&mxproject.PreviousLibFiles.LibFiles, section, "LibFiles")
	}

	section = inidata.Section("PreviousUsedKeilFiles")
	if section != nil {
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.SourceFiles, section, "SourceFiles")
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.HeaderPath, section, "HeaderPath")
		StoreItemCsv(&mxproject.PreviousUsedKeilFiles.CDefines, section, "CDefines")

		PrintItemCsv(section, "SourceFiles")
		PrintItemCsv(section, "HeaderPath")
		PrintItemCsv(section, "CDefines")
	}

	section = inidata.Section("PreviousGenFiles")
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
