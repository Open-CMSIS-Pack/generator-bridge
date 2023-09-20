/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"fmt"
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

func WriteYml() {
	fruits := [...]string{"apple", "orange", "mango", "strawberry"}
	data, err := yaml.Marshal(fruits)
	if err != nil {
		log.Fatal(err)
	}
	err1 := os.WriteFile("fruits.yaml", data, 0644)
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("Success!")
}
