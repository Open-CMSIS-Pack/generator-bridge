/*
 * Copyright (c) 2023-2026 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func getCgenLogPath(cgenPath string) string {
	return strings.TrimSuffix(cgenPath, ".yml") + ".log"
}

// deleteCgenLog deletes the *.cgen.log file for a given cgen path.
func deleteCgenLog(cgenPath string) error {
	logPath := getCgenLogPath(cgenPath)
	if utils.FileExists(logPath) {
		return os.Remove(logPath)
	}
	return nil
}

// logCgenError appends an error message with timestamp and function name to the *.cgen.log file.
func logCgenError(cgenPath string, err error) {
	logPath := getCgenLogPath(cgenPath)

	// Get caller's function name.
	pc, _, _, _ := runtime.Caller(1)
	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		name := fn.Name()
		if idx := strings.LastIndex(name, "/"); idx >= 0 && idx+1 < len(name) {
			name = name[idx+1:]
		}
		funcName = name
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf("[%s] %s: %v\n", timestamp, funcName, err)

	file, openErr := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openErr != nil {
		log.Warnf("failed to open cgen log '%s': %v", logPath, openErr)
		return
	}
	defer file.Close()

	_, writeErr := file.WriteString(message)
	if writeErr != nil {
		log.Warnf("failed to write cgen log '%s': %v", logPath, writeErr)
	}
}

// logCgenInfoIfLogExists appends an info message only when the cgen log already
// exists. This is used in daemon mode to acknowledge recovery after prior
// errors/warnings in the same session.
func logCgenInfoIfLogExists(cgenPath, message string) {
	logPath := getCgenLogPath(cgenPath)
	if !utils.FileExists(logPath) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	entry := fmt.Sprintf("[%s] info: %s\n", timestamp, message)

	file, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if openErr != nil {
		log.Warnf("failed to open cgen log '%s': %v", logPath, openErr)
		return
	}
	defer file.Close()

	if _, writeErr := file.WriteString(entry); writeErr != nil {
		log.Warnf("failed to write cgen log '%s': %v", logPath, writeErr)
	}
}

// deleteAllCgenLogs deletes all *.cgen.log files for the given bridge parameters.
func deleteAllCgenLogs(bridgeParams []BridgeParamType) {
	for _, bp := range bridgeParams {
		_ = deleteCgenLog(bp.CgenName)
	}
}
