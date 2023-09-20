/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32CubeMX

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

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

func IniReader(path string) error {
	log.Infof("\nReading CubeMX config file: %v", path)

	inidata, err := ini.Load(path)
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return nil
	}

	section := inidata.Section("PreviousUsedKeilFiles")
	if section != nil {
		key := "SourceFiles"
		valStr := section.Key(key).Strings(";")
		PrintKeyValStrs(key, valStr)

		key = "HeaderPath"
		valStr = section.Key(key).Strings(";")
		PrintKeyValStrs(key, valStr)

		key = "CDefines"
		valStr = section.Key(key).Strings(";")
		PrintKeyValStrs(key, valStr)
	}

	return nil
}

/*
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
*/
