/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"errors"
	"fmt"
	"io"
	"strings"

	readfile "github.com/open-cmsis-pack/generator-bridge/internal/readFile"
	stm32cubemx "github.com/open-cmsis-pack/generator-bridge/internal/stm32CubeMX"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AllCommands contains all available commands for generator-bridge
var AllCommands = []*cobra.Command{}

var rootCommand *cobra.Command = nil

func GetConfig() *pflag.FlagSet {
	return rootCommand.PersistentFlags()
}

// configureInstaller configures generator-bridge installer for adding or removing pack/pdsc
func configureGlobalCmd(cmd *cobra.Command, args []string) error {
	verbosiness, _ := GetConfig().GetBool("verbose")
	quiet, _ := GetConfig().GetBool("quiet")
	if quiet && verbosiness {
		return errors.New("both \"-q\" and \"-v\" were specified, please pick only one verboseness option")
	}

	log.SetLevel(log.InfoLevel)
	log.SetOutput(cmd.OutOrStdout())

	if quiet {
		log.SetLevel(log.ErrorLevel)
	}

	if verbosiness {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

var flags struct {
	version bool
	help    bool
	inFile  string
	outPath string
}

var Version string
var CopyRight string

func printVersionAndLicense(file io.Writer) {
	fmt.Fprintf(file, "generator-bridge version %v %s\n", strings.ReplaceAll(Version, "v", ""), CopyRight)
}

// UsageTemplate returns usage template for the command.
var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func NewCli() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "generator-bridge [command] [flags]",
		Short:             "This utility is a bridge to Vendor tools, e.g. STCube",
		Long:              "Please refer to the upstream repository for further information: https://github.com/Open-CMSIS-Pack/generator-bridge.",
		SilenceUsage:      true,
		SilenceErrors:     true,
		PersistentPreRunE: configureGlobalCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.version {
				printVersionAndLicense(cmd.OutOrStdout())
				return nil
			}

			if flags.help {
				return cmd.Help()
			}

			if flags.inFile != "" {
				return readfile.Process(flags.inFile, flags.outPath)
			}

			if len(args) == 1 {
				cbuildYmlPath := args[0]
				return stm32cubemx.Process(cbuildYmlPath, flags.outPath, "")
			}

			return cmd.Help()
		},
	}

	rootCmd.SetUsageTemplate(usageTemplate)

	rootCmd.Flags().BoolVarP(&flags.version, "version", "V", false, "Prints the version number of generator-bridge and exit")
	rootCmd.Flags().BoolVarP(&flags.help, "help", "h", false, "Show help")
	rootCmd.Flags().StringVarP(&flags.inFile, "read", "r", "", "Reads an input file, type is auto determined")
	rootCmd.Flags().StringVarP(&flags.outPath, "out", "o", "", "Output path for generated files")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run silently, printing only error messages")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Sets verboseness level: None (Errors + Info + Warnings), -v (all + Debugging). Specify \"-q\" for no messages")

	for _, cmd := range AllCommands {
		rootCmd.AddCommand(cmd)
	}

	rootCommand = rootCmd

	return rootCmd
}
