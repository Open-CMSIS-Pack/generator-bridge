/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type TextBuilder struct {
	text string
}

func AddQuotes(text string) string {
	return "\"" + text + "\""
}

func (t *TextBuilder) AddLine(args ...string) {
	for _, arg := range args {
		if t.text != "" && t.text[len(t.text)-1] != '\n' {
			t.text += " "
		}
		t.text += arg
	}
	t.text += "\n"
}

func (t *TextBuilder) AddSpaces(num, tabWidth int) {
	for i := int(0); i < num; i++ {
		for j := int(0); j < tabWidth; j++ {
			t.text += " "
		}
	}
}

func (t *TextBuilder) GetLine() string {
	return t.text
}

// FileExists checks if filePath is an actual file in the local file system
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if dirPath is an actual directory in the local file system
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// EnsureDir recursevily creates a directory tree if it doesn't exist already
func EnsureDir(dirName string) error {
	log.Debugf("Ensuring \"%s\" directory exists", dirName)
	err := os.MkdirAll(dirName, 0755)
	if err != nil && !os.IsExist(err) {
		log.Error(err)
		return nil //errs.ErrFailedCreatingDirectory
	}
	return nil
}

func GetDateTimeString() string {
	currentTime := time.Now()
	text := currentTime.Format("2006-01-02 15:04:05")

	return text
}
