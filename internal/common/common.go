/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"bytes"
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

func WriteYml(path string, out interface{}) error {
	var data bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&data)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&out)
	if err != nil {
		return err
	}

	err1 := os.WriteFile(path, data.Bytes(), 0600)
	if err1 != nil {
		log.Fatal(err1)
	}

	return nil
}
