/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/ini.v1"
)

func ExamplePrintKeyValStr() {
	PrintKeyValStr("key", "val")
	// Output:
	//
	// key : val
}

func ExamplePrintKeyValStrs() {
	PrintKeyValStrs("key", []string{"val1", "val2"})
	// Output:
	//
	// key
	// 0: val1
	// 1: val2
}

func ExamplePrintKeyValInt() {
	PrintKeyValInt("key", 4711)
	// Output:
	//
	// key : 4711
}

// Test_GetData verifies parsing of different sections (.mxproject logic)
func Test_GetData(t *testing.T) {
	t.Parallel()

	// Helper function to write a temporary INI file
	writeIni := func(t *testing.T, content string) *ini.File {
		t.Helper()
		tmp := filepath.Join(t.TempDir(), "test.mxproject")
		if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
			t.Fatalf("failed writing tmp ini: %v", err)
		}
		defer func() {
			_ = os.Remove(tmp)
		}()
		inidata, err := ini.Load(tmp)
		if err != nil {
			t.Fatalf("failed loading ini: %v", err)
		}
		return inidata
	}

	iniFull := `
[Ctx1:ThirdPartyIp]
ThirdPartyIpNumber=2
ThirdPartyIpName#0=IPX
ThirdPartyIpName#1=IPY

[Ctx1:ThirdPartyIp#IPX]
include=inc1 ;inc2;inc3
sourceAsm=asm1 ;asm2
source=ipx_src1.c ;ipx_src2.c

[Ctx1:ThirdPartyIp#IPY]
include=incY1 ;incY2
source=ipy_src.c ;ipy_src_more.c

[Ctx1:PreviousUsedCubeIDEFiles]
SourceFiles=src1.c ;src2.c;src3.c
HeaderPath=./include ;./inc2;./inc3
CDefines=DEF1 ;DEF_TWO;DEF_THREE

[Ctx1:PreviousLibFiles]
LibFiles=libA.a ;libB.a

[Ctx1:PreviousGenFiles]
AdvancedFolderStructure=1
HeaderFileListSize=2
HeaderFiles#0=hf0.h
HeaderFiles#1=hf1.h
HeaderFolderListSize=1
HeaderPath#0=./gen/inc
HeaderFiles=headers.h
SourceFileListSize=1
SourceFiles#0=sf0.c
SourceFolderListSize=1
SourcePath#0=./gen/src
SourceFiles=sfMain.c
`

	iniNoCtx := `
[ThirdPartyIp]
ThirdPartyIpNumber=1
ThirdPartyIpName#0=LIB

[ThirdPartyIp#LIB]
include=lib_inc ;lib_extra
sourceAsm=lib_startup.s
source=lib_src.c ;lib_more.c

[PreviousUsedKeilFiles]
SourceFiles=main.c ;utils.c
HeaderPath=./inc ;./inc2
CDefines=KEIL_DEF ;OTHER
`

	tests := []struct {
		name     string
		inidata  *ini.File
		ctx      string
		compiler string
		wantErr  bool
		check    func(t *testing.T, mx MxprojectType)
	}{
		{
			name:     "context_gcc_full_sections",
			inidata:  writeIni(t, iniFull),
			ctx:      "Ctx1",
			compiler: "GCC",
			wantErr:  false,
			check: func(t *testing.T, mx MxprojectType) {
				// ThirdPartyIpFiles structure changed: now slice of ThirdPartyIpNames
				var includes, asms, sources []string
				var names []string
				for _, ip := range mx.ThirdPartyIpFiles {
					includes = append(includes, ip.IncludeFiles...)
					asms = append(asms, ip.SourceAsmFiles...)
					sources = append(sources, ip.SourceFiles...)
					names = append(names, ip.ThirdPartyIpName)
				}
				if !containsAll(names, []string{"IPX", "IPY"}) {
					t.Errorf("expected IP names IPX/IPY got %v", names)
				}
				if !containsAll(includes, []string{"inc1", "inc2", "inc3", "incY1", "incY2"}) {
					t.Errorf("expected all include files, got: %v", includes)
				}
				if !containsAll(asms, []string{"asm1", "asm2"}) { // IPY has no sourceAsm
					t.Errorf("expected asm1/asm2 in SourceAsmFiles: %v", asms)
				}
				if !containsAll(sources, []string{"ipx_src1.c", "ipx_src2.c", "ipy_src.c", "ipy_src_more.c"}) {
					t.Errorf("third-party source files incomplete: %v", sources)
				}
				// PreviousUsedFiles
				expSrc := []string{"src1.c", "src2.c", "src3.c"}
				if !reflect.DeepEqual(mx.PreviousUsedFiles.SourceFiles, expSrc) {
					t.Errorf("SourceFiles expected %v got %v", expSrc, mx.PreviousUsedFiles.SourceFiles)
				}
				expInc := []string{"./include", "./inc2", "./inc3"}
				if !reflect.DeepEqual(mx.PreviousUsedFiles.HeaderPath, expInc) {
					t.Errorf("HeaderPath expected %v got %v", expInc, mx.PreviousUsedFiles.HeaderPath)
				}
				expDef := []string{"DEF1", "DEF_TWO", "DEF_THREE"}
				if !reflect.DeepEqual(mx.PreviousUsedFiles.CDefines, expDef) {
					t.Errorf("CDefines expected %v got %v", expDef, mx.PreviousUsedFiles.CDefines)
				}
				// PreviousLibFiles
				if !reflect.DeepEqual(mx.PreviousLibFiles.LibFiles, []string{"libA.a", "libB.a"}) {
					t.Errorf("LibFiles mismatch: %v", mx.PreviousLibFiles.LibFiles)
				}
				// PreviousGenFiles (only sample check)
				if mx.PreviousGenFiles.AdvancedFolderStructure != "1" {
					t.Errorf("AdvancedFolderStructure expected '1' got %v", mx.PreviousGenFiles.AdvancedFolderStructure)
				}
				if !containsAll(mx.PreviousGenFiles.HeaderFilesList, []string{"hf0.h", "hf1.h"}) {
					t.Errorf("HeaderFilesList missing entries: %v", mx.PreviousGenFiles.HeaderFilesList)
				}
				if !containsAll(mx.PreviousGenFiles.SourceFilesList, []string{"sf0.c"}) {
					t.Errorf("SourceFilesList expected sf0.c: %v", mx.PreviousGenFiles.SourceFilesList)
				}
				if mx.PreviousGenFiles.SourceFiles != "sfMain.c" {
					t.Errorf("single SourceFiles expected sfMain.c got %v", mx.PreviousGenFiles.SourceFiles)
				}
			},
		},
		{
			name:     "no_context_keil",
			inidata:  writeIni(t, iniNoCtx),
			ctx:      "",
			compiler: "AC6",
			wantErr:  false,
			check: func(t *testing.T, mx MxprojectType) {
				var includes []string
				for _, ip := range mx.ThirdPartyIpFiles {
					includes = append(includes, ip.IncludeFiles...)
				}
				if !containsAll(includes, []string{"lib_inc", "lib_extra"}) {
					t.Errorf("IncludeFiles expected lib_inc/lib_extra: %v", includes)
				}
				if !containsAll(mx.PreviousUsedFiles.SourceFiles, []string{"main.c", "utils.c"}) {
					t.Errorf("SourceFiles expected main.c/utils.c: %v", mx.PreviousUsedFiles.SourceFiles)
				}
				if !containsAll(mx.PreviousUsedFiles.CDefines, []string{"KEIL_DEF", "OTHER"}) {
					t.Errorf("CDefines expected KEIL_DEF/OTHER: %v", mx.PreviousUsedFiles.CDefines)
				}
			},
		},
		{
			name:     "unknown_compiler_error",
			inidata:  writeIni(t, iniFull),
			ctx:      "Ctx1",
			compiler: "UNKNOWN",
			wantErr:  true,
			check:    func(t *testing.T, mx MxprojectType) {},
		},
		{
			name:     "empty_ini_no_sections",
			inidata:  writeIni(t, ""),
			ctx:      "Ctx1",
			compiler: "GCC",
			wantErr:  false,
			check: func(t *testing.T, mx MxprojectType) {
				if len(mx.PreviousUsedFiles.SourceFiles) != 0 || len(mx.ThirdPartyIpFiles) != 0 {
					t.Errorf("expected empty arrays/slices for missing sections, got %+v", mx)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mx, err := GetData(tt.inidata, tt.ctx, tt.compiler)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetData() error=%v wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			tt.check(t, mx)
		})
	}
}

// containsAll checks whether all expected values are present in the slice
func containsAll(have []string, expected []string) bool {
	set := make(map[string]struct{}, len(have))
	for _, v := range have {
		set[v] = struct{}{}
	}
	for _, e := range expected {
		if _, ok := set[e]; !ok {
			return false
		}
	}
	return true
}
