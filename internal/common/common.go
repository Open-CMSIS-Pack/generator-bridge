/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"bytes"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ReadYml(path string, out interface{}) error {

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("yamlFile.Get err %v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
		return err
	}

	return nil
}

// WriteYml serializes out as YAML and writes it to path only when the on-disk content differs.
// It returns true if a write occurred, or false if the existing file was already up-to-date.
// Callers must not concurrently modify the destination file.
func WriteYml(path string, out interface{}) (bool, error) {
	var data bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&data)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(&out); err != nil {
		return false, fmt.Errorf("encode YAML %q: %w", path, err)
	}
	if err := yamlEncoder.Close(); err != nil {
		return false, fmt.Errorf("close YAML encoder for %q: %w", path, err)
	}

	// If the existing file cannot be read, attempt the write because content
	// comparison is an optimization and writing the generated YAML is primary.
	existingData, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existingData, data.Bytes()) {
		return false, nil
	}

	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return false, fmt.Errorf("write YAML %q: %w", path, err)
	}

	return true, nil
}
