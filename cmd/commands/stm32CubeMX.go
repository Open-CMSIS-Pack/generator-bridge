/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package commands

import (
	"github.com/open-cmsis-pack/generator-bridge/cmd/stm32CubeMX"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stm32CubeMXCmdFlags struct {
	// Reports encoded progress for files and download when used by other tools
	path string
	args *[]string
}

var STM32CubeMXCmd = &cobra.Command{
	Use:               "STM32CubeMX",
	Short:             "Launch STM32CubeMX",
	Long:              getLongUpdateDescription(),
	PersistentPreRunE: nil,
	Args:              cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Launch STM32CubeMX")
		err := stm32CubeMX.Launch(stm32CubeMXCmdFlags.path, stm32CubeMXCmdFlags.args)

		return err
	},
}

func getLongUpdateDescription() string {
	return `Launch STM32CubeMX Command`
}

func init() {
	STM32CubeMXCmd.Flags().StringVar(&stm32CubeMXCmdFlags.path, "launch", "", "Cube-MX input file .mxproject")
	stm32CubeMXCmdFlags.args = STM32CubeMXCmd.Flags().StringArray("args", []string{}, "Arguments for CubeMX launch")
}
