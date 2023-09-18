/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package stm32CubeMX

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Launch(path string, args *[]string) error {
	log.Infof("Launching STM32CubeMX...")

	cmd := exec.Command(path, *args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	IniReader("")

	return nil
}
