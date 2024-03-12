/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"os"
	"time"

	"github.com/open-cmsis-pack/generator-bridge/cmd/commands"
	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(LogFormatter))
	f, err := os.OpenFile("cbridge.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		defer f.Close()
		log.SetOutput(f)
	}

	utils.StartSignalWatcher()
	start := time.Now()

	log.Println("Command line:", os.Args[1:])

	commands.Version = version
	commands.Copyright = copyright
	cmd := commands.NewCli()
	err = cmd.Execute()
	if err != nil {
		log.Errorf("Error : %v", err)
		os.Exit(-1)
	}

	log.Debugf("Took %v", time.Since(start))
	utils.StopSignalWatcher()
}
