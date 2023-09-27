/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Copy of cmd/log.go
type LogFormatter struct{}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	msg := fmt.Sprintf("%s: %s\n", level[0:1], entry.Message)
	return []byte(msg), nil
}

var testingDir = filepath.Join("..", "..", "testdata", "integration")

type TestCase struct {
	args           []string
	name           string
	defaultMode    bool
	createPackRoot bool
	expectedStdout []string
	expectedStderr []string
	expectedErr    error
	setUpFunc      func(t *TestCase)
	tearDownFunc   func()
	validationFunc func(t *testing.T)
	assert         *assert.Assertions
	env            map[string]string
}

func init() {
	logLevel := log.InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(new(LogFormatter))
}
