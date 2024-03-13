/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func Test_configureGlobalCmd(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test", args{cmd: &cobra.Command{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewCli()
			if err := configureGlobalCmd(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("configureGlobalCmd() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func Test_printVersionAndLicense(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		copyright string
		wantFile  string
	}{
		{"test", "v1.2.3", "copy", "generator-bridge version 1.2.3 copy\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &bytes.Buffer{}
			saveVersion := Version
			saveCopyright := Copyright
			Version = tt.version
			Copyright = tt.copyright
			printVersionAndLicense(file)
			if gotFile := file.String(); gotFile != tt.wantFile {
				t.Errorf("printVersionAndLicense() %s = %v, want %v", tt.name, gotFile, tt.wantFile)
			}
			Version = saveVersion
			Copyright = saveCopyright
		})
	}
}
