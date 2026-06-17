/*
 * Copyright (c) 2023-2026 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
)

func Test_logCgenError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cgenPath := filepath.Join(tmpDir, "test.cgen.yml")
	logPath := strings.TrimSuffix(cgenPath, ".yml") + ".log"

	// Test writing an error
	logCgenError(cgenPath, errors.New("test error 1"))

	// Verify the log file exists
	if !utils.FileExists(logPath) {
		t.Errorf("expected log file to be created at %s", logPath)
	}

	// Verify the error message was written with timestamp
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "test error 1") {
		t.Errorf("expected error message in log file, got: %s", content)
	}
	if !strings.Contains(content, "Test_logCgenError") {
		t.Errorf("expected function name in log file, got: %s", content)
	}

	// Test appending a second error
	logCgenError(cgenPath, errors.New("test error 2"))

	data, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file after append: %v", err)
	}

	content = string(data)
	if !strings.Contains(content, "test error 1") || !strings.Contains(content, "test error 2") {
		t.Errorf("expected both error messages in log file, got: %s", content)
	}
	if !strings.Contains(content, "Test_logCgenError") {
		t.Errorf("expected function name in log file, got: %s", content)
	}
}

func Test_deleteCgenLog(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cgenPath := filepath.Join(tmpDir, "test.cgen.yml")
	logPath := strings.TrimSuffix(cgenPath, ".yml") + ".log"

	// Create a log file
	logCgenError(cgenPath, errors.New("test error"))

	if !utils.FileExists(logPath) {
		t.Fatalf("log file was not created")
	}

	// Delete the log file
	err := deleteCgenLog(cgenPath)
	if err != nil {
		t.Fatalf("deleteCgenLog() returned unexpected error: %v", err)
	}

	// Verify the log file is gone
	if utils.FileExists(logPath) {
		t.Errorf("expected log file to be deleted, but it still exists at %s", logPath)
	}

	// Test deleting a non-existent log file (should not error)
	err = deleteCgenLog(cgenPath)
	if err != nil {
		t.Fatalf("deleteCgenLog() should not error for non-existent file: %v", err)
	}
}

func Test_deleteAllCgenLogs(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Create multiple cgen.log files
	cgenPath1 := filepath.Join(tmpDir, "cgen1.yml")
	cgenPath2 := filepath.Join(tmpDir, "cgen2.yml")

	logCgenError(cgenPath1, errors.New("error 1"))
	logCgenError(cgenPath2, errors.New("error 2"))

	logPath1 := strings.TrimSuffix(cgenPath1, ".yml") + ".log"
	logPath2 := strings.TrimSuffix(cgenPath2, ".yml") + ".log"

	if !utils.FileExists(logPath1) || !utils.FileExists(logPath2) {
		t.Fatalf("log files were not created")
	}

	// Delete all logs
	bridgeParams := []BridgeParamType{
		{CgenName: cgenPath1},
		{CgenName: cgenPath2},
	}
	deleteAllCgenLogs(bridgeParams)

	// Verify both log files are gone
	if utils.FileExists(logPath1) {
		t.Errorf("expected log file 1 to be deleted")
	}
	if utils.FileExists(logPath2) {
		t.Errorf("expected log file 2 to be deleted")
	}
}
