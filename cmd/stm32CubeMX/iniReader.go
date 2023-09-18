/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32CubeMX

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func PrintKeyValStr(key, val string) {
	fmt.Printf("\n%v : %v", key, val)
}

func PrintKeyValInt(key string, val int) {
	fmt.Printf("\n%v : %v", key, val)
}

func IniReader(path string) error {
	log.Infof("\nReading CubeMX config file: %v", path)

	inidata, err := ini.Load(path)
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}
	section := inidata.Section("CortexM33S:PreviousGenFiles")

	key := "AdvancedFolderStructure"
	valStr := section.Key(key).String()
	PrintKeyValStr(key, valStr)

	key = "HeaderFileListSize"
	valInt, err := section.Key(key).Int()
	if err == nil {
		PrintKeyValInt(key, valInt)
	}

	cnt := 0
	key = "HeaderFiles#"
	for {
		keyN := key + strconv.Itoa(cnt)
		valStr = section.Key(keyN).String()
		if valStr == "" {
			break
		}
		fmt.Printf("\n%v : %v", keyN, valStr)
		cnt++
	}

	return nil
}
