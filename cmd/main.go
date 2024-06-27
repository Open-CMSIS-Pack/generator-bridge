/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"io"
	"os"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/cmd/commands"
	stm32cubemx "github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(LogFormatter))

	utils.StartSignalWatcher()
	start := time.Now()

	commands.Version = version
	commands.Copyright = copyright
	cmd := commands.NewCli()
	err := cmd.Execute()
	if err != nil {
		log.Errorf("Error : %v", err)
		os.Exit(-1)
	}

	log.Debugf("Took %v", time.Since(start))
	if stm32cubemx.LogFile != nil {
		_ = stm32cubemx.LogFile.Close()
	}
	log.SetOutput((io.Discard))
	utils.StopSignalWatcher()
}
