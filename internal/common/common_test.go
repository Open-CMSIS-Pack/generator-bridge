/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestReadYml(t *testing.T) {
	type TestYml struct {
		Test []struct {
			Xx string `yaml:"xx"`
		} `yaml:"test"`
	}
	var testyml TestYml

	t1 := TestYml{Test: []struct {
		Xx string `yaml:"xx"`
	}{{Xx: "abc"}}}

	type args struct {
		path string
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    TestYml
		wantErr bool
	}{
		{"test", args{"../../testdata/test.yml", &testyml}, t1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadYml(tt.args.path, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("ReadYml() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			xx := tt.want
			if !reflect.DeepEqual(tt.args.out, &xx) {
				t.Errorf("ReadYml() %s got = %v, want %v", tt.name, tt.args.out, tt.want)
			}
		})
	}
}

func TestWriteYmlDoesNotTouchUnchangedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.cgen.yml")
	content := struct {
		Value string `yaml:"value"`
	}{Value: "unchanged"}

	changed, err := WriteYml(path, &content)
	if err != nil {
		t.Fatalf("initial WriteYml() error = %v", err)
	}
	if !changed {
		t.Fatal("initial WriteYml() changed = false, want true")
	}

	modTime := time.Date(2020, time.January, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("os.Chtimes() error = %v", err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat() before second write error = %v", err)
	}

	changed, err = WriteYml(path, &content)
	if err != nil {
		t.Fatalf("second WriteYml() error = %v", err)
	}
	if changed {
		t.Fatal("second WriteYml() changed = true, want false")
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat() after second write error = %v", err)
	}
	if !after.ModTime().Equal(before.ModTime()) {
		t.Errorf("modification time changed from %v to %v", before.ModTime(), after.ModTime())
	}
}

func TestWriteYmlReturnsWriteError(t *testing.T) {
	path := t.TempDir()
	content := struct {
		Value string `yaml:"value"`
	}{Value: "test"}

	changed, err := WriteYml(path, &content)
	if err == nil {
		t.Fatal("WriteYml() error = nil, want write error")
	}
	if changed {
		t.Fatal("WriteYml() changed = true after failed write, want false")
	}
}
