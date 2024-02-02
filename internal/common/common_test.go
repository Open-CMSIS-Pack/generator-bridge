/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package common

import (
	"reflect"
	"testing"
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
			if !reflect.DeepEqual(tt.args.out, &tt.want) {
				t.Errorf("ReadYml() %s got = %v, want %v", tt.name, tt.args.out, tt.want)
			}
		})
	}
}
