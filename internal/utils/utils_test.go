/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"os"
	"path"
	"testing"

	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	"github.com/stretchr/testify/assert"
)

var testDir string = "./Testing"

func TestAddLine(t *testing.T) {
	var text utils.TextBuilder

	text.AddLine("A line")
	expected := "A line\n"
	result := text.GetLine()
	assert.Equal(t, expected, result)

	text.AddLine("A second line")
	result = text.GetLine()
	expected += "A second line\n"
	assert.Equal(t, expected, result)
}

func TestAddQuotes(t *testing.T) {
	text := "Test"
	expected := "\"" + text + "\""
	result := utils.AddQuotes(text)
	assert.Equal(t, expected, result)
}

func TestFileExists(t *testing.T) {

	result := utils.DirExists(testDir)
	expected := true
	assert.Equal(t, expected, result)

	filename := path.Join(testDir, "fileexists.txt")
	result = utils.FileExists(filename)
	expected = false
	assert.Equal(t, expected, result)

	text := "Hello, World"
	os.WriteFile(filename, []byte(text), 0755)
	result = utils.FileExists(filename)
	expected = false
	assert.Equal(t, expected, result)
}

func Init() {
	os.MkdirAll(testDir, 0755)

}
