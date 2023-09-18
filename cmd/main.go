/*
 * Copyright (c) 2022-2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"os"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/cmd/commands"
	"github.com/open-cmsis-pack/generator-bridge/cmd/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(LogFormatter))
	log.SetOutput(os.Stdout)

	utils.StartSignalWatcher()
	start := time.Now()

	commands.Version = version
	commands.CopyRight = copyRight
	cmd := commands.NewCli()
	err := cmd.Execute()
	if err != nil {
		os.Exit(-1)
	}

	log.Debugf("Took %v", time.Since(start))
	utils.StopSignalWatcher()
}
