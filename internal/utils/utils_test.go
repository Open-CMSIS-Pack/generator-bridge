/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestAddQuotes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{"Test"}, "\"Test\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddQuotes(tt.args.text); got != tt.want {
				t.Errorf("AddQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTextBuilder_AddLine(t *testing.T) {
	var builder TextBuilder
	var builder1 TextBuilder
	var builder2 TextBuilder

	type args struct {
		args []string
	}
	tests := []struct {
		name string
		tr   *TextBuilder
		args args
		want string
	}{
		{"nix", &builder, args{}, "\n"},
		{"1", &builder1, args{[]string{"A line"}}, "A line\n"},
		{"1+", &builder2, args{[]string{"A line", "and more"}}, "A line and more\n"},
		{"2", &builder1, args{[]string{"A second line"}}, "A line\nA second line\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.AddLine(tt.args.args...)
			if !reflect.DeepEqual(tt.tr.GetLine(), tt.want) {
				t.Errorf("TestTextBuilder_AddLine() %s got = %v, want %v", tt.name, tt.tr.GetLine(), tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	t.Parallel()
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"file", args{"../../testdata/test.yml"}, true},
		{"dir", args{"../../testdata"}, false},
		{"nix", args{"../../testdata/nix"}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := FileExists(tt.args.filePath); got != tt.want {
				t.Errorf("FileExists() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	type args struct {
		dirPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"file", args{"../../testdata/test.yml"}, false},
		{"dir", args{"../../testdata"}, true},
		{"nix", args{"../../testdata/nix"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirExists(tt.args.dirPath); got != tt.want {
				t.Errorf("DirExists() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	type args struct {
		dirName string
	}
	tests := []struct {
		name    string
		args    args
		remove  string
		wantErr bool
	}{
		{"test", args{"../../testdata/1/2/3"}, "../../testdata/1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.remove)
			if err := EnsureDir(tt.args.dirName); (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestConvertFilename(t *testing.T) {
	type args struct {
		outPath         string
		file            string
		relativePathAdd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"testAbs", args{"../../testdata", "C:/test.ioc", "stm32cubemx"}, "C:/test.ioc", false},
		{"test", args{"../../testdata", "test.ioc", "stm32cubemx"}, "./stm32cubemx/test.ioc", false},
		{"nix", args{"../../testdata", "nix", "stm32cubemx"}, "./stm32cubemx/nix", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !(tt.name == "testAbs" && runtime.GOOS == "windows") {
				got, err := ConvertFilename(tt.args.outPath, tt.args.file, tt.args.relativePathAdd)
				if (err != nil) != tt.wantErr {
					t.Errorf("ConvertFilename() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ConvertFilename() %s = %v, want %v", tt.name, got, tt.want)
				}
			}
		})
	}
}
