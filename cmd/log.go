/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// LogFormatter is generator-bridge's basic log formatter
type LogFormatter struct{}

// Format prints out logs like "I: some message", where the first letter indicates (I)NFO, (D)EBUG, (W)ARNING or (E)RROR
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	// level := strings.ToUpper(entry.Level.String())
	// msg := fmt.Sprintf("%s: %s\n", level[0:1], entry.Message)
	msg := fmt.Sprintf("%s\n", entry.Message)
	return []byte(msg), nil
}
