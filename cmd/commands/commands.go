/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the generator-bridge project. */

package commands

import (
	"errors"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	viperType "github.com/spf13/viper"
)

// AllCommands contains all available commands for generator-bridge
var AllCommands = []*cobra.Command{
	STM32CubeMXCmd,
}

var viper *viperType.Viper

func configureInstallerGlobalCmd(cmd *cobra.Command, args []string) error {
	verbosiness := viper.GetBool("verbose")
	quiet := viper.GetBool("quiet")
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

// configureInstaller configures generator-bridge installer for adding or removing pack/pdsc
func configureInstaller(cmd *cobra.Command, args []string) error {
	err := configureInstallerGlobalCmd(cmd, args)
	if err != nil {
		return err
	}

	return nil
}

var flags struct {
	version bool
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
		Use:           "generator-bridge [command] [flags]",
		Short:         "This utility is a bridge to Vendor tools, e.g. STCube",
		Long:          "Please refer to the upstream repository for further information: https://github.com/Open-CMSIS-Pack/generator-bridge.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.version {
				printVersionAndLicense(cmd.OutOrStdout())
				return nil
			}

			return cmd.Help()
		},
	}

	rootCmd.SetUsageTemplate(usageTemplate)

	viper = viperType.New()

	rootCmd.Flags().BoolVarP(&flags.version, "version", "V", false, "Prints the version number of generator-bridge and exit")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Run generator-bridge silently, printing only error messages")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Sets verboseness level: None (Errors + Info + Warnings), -v (all + Debugging). Specify \"-q\" for no messages")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	for _, cmd := range AllCommands {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}
