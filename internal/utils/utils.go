/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"errors"
	"os"
	"path/filepath"

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

func ConvertFilename(outPath, file, relativePathAdd string) (string, error) {
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	// Check if the file is absolute
	if !filepath.IsAbs(file) {
		toolchainPath := filepath.Join(outPath, relativePathAdd) // create the path where STCube sets it's files relative to toolchain folder( example :./STM32CubeMX/MDK-ARM/)
		file = filepath.Join(toolchainPath, file)
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Errorf("file or directory not found: %v", file)
	}

	var err error
	origfilename := file
	file, err = filepath.Rel(outPath, file)
	if err != nil {
		log.Errorf("path error found: %v", file)
		return origfilename, nil
	}
	file = filepath.ToSlash(file)
	file = "./" + file

	return file, nil
}

func ConvertFilenameRel(outPath, file string) (string, error) {
	file = filepath.Clean(file)
	file = filepath.ToSlash(file)

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Errorf("file or directory not found: %v", file)
	}

	var err error
	origfilename := file
	file, err = filepath.Rel(outPath, file)
	if err != nil {
		log.Errorf("path error found: %v", file)
		return origfilename, errors.New("path error")
	}

	file = filepath.ToSlash(file)
	file = "./" + file

	return file, nil
}
