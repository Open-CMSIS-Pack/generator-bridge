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
	path struct {
		cubeMx    string
		cbuildYml string
	}
	args *[]string
}

var STM32CubeMXCmd = &cobra.Command{
	Use:               "STM32CubeMX",
	Short:             "Launch STM32CubeMX",
	Long:              `Launch STM32CubeMX Command`,
	PersistentPreRunE: nil,
	Args:              cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Launch STM32CubeMX")
		err := stm32CubeMX.Process(stm32CubeMXCmdFlags.path.cbuildYml, stm32CubeMXCmdFlags.path.cubeMx)

		return err
	},
}

func init() {
	STM32CubeMXCmd.Flags().StringVar(&stm32CubeMXCmdFlags.path.cubeMx, "launch", "", "Cube-MX input file .mxproject")
	STM32CubeMXCmd.Flags().StringVar(&stm32CubeMXCmdFlags.path.cbuildYml, "cbuildYml", "", "CBuild YAML input file")
	stm32CubeMXCmdFlags.args = STM32CubeMXCmd.Flags().StringArray("args", []string{}, "Arguments for CubeMX launch")
}
