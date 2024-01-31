/*
 * Copyright (c) 2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package generator

import (
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	var params ParamsType

	type args struct {
		name   string
		params *ParamsType
	}
	tests := []struct {
		name    string
		args    args
		want    ParamsType
		wantErr bool
	}{
		{"wrong.yml", args{"../../testdata/wrong.yml", &params}, ParamsType{}, true},
		{"global.yml", args{"../../testdata/global.yml", &params}, ParamsType{ID: "CubeMX", DownloadURL: "https://nix.html"}, false},
		{"global-nix.yml", args{"../../testdata/global-nix.yml", &params}, ParamsType{}, true},
		{"wrong.yml", args{"../../testdata/wrong.yml", &params}, ParamsType{}, true},
		{"xxx.yml", args{"../../testdata/xxx.yml", &params}, ParamsType{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params.ID = ""
			params.DownloadURL = ""
			if err := Read(tt.args.name, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Read() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(params, tt.want) {
				t.Errorf("createContextMap() %s got = %v, want %v", tt.name, params, tt.want)
			}
		})
	}
}
