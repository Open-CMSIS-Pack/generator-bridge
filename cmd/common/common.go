/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ReadYml(path string, out interface{}) error {

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Infof("yamlFile.Get err %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
	}

	return nil
}

func WriteYml(path, header string, out interface{}) error {
	data, err := yaml.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	err1 := os.WriteFile(path, data, 0644)
	if err1 != nil {
		log.Fatal(err1)
	}

	return nil
}
