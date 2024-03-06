/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
)

func Test_GetToolchain(t *testing.T) {
	t.Parallel()

	type args struct {
		compiler string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"AC6"}, "MDK-ARM V5", false},
		{"test2", args{"GCC"}, "STM32CubeIDE", false},
		{"test3", args{"IAR"}, "EWARM", false},
		{"test4", args{"CLANG"}, "STM32CubeIDE", false},
		{"fail", args{"unknown"}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetToolchain(tt.args.compiler)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToolchain() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetToolchain() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_GetRelativePathAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		outPath  string
		compiler string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"./STM32CubeMX", "AC6"}, "MDK-ARM", false},
		{"test2", args{"./STM32CubeMX", "GCC"}, "", false},
		{"test3", args{"./STM32CubeMX", "IAR"}, "EWARM", false},
		{"test4", args{"./STM32CubeMX", "CLANG"}, "", false},
		{"test5", args{"./", "AC6"}, "STM32CubeMX/MDK-ARM", false},
		{"test6", args{"./", "GCC"}, "STM32CubeMX", false},
		{"test7", args{"./", "IAR"}, "STM32CubeMX/EWARM", false},
		{"test8", args{"./", "CLANG"}, "STM32CubeMX", false},
		{"fail", args{"", "unknown"}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetRelativePathAdd(tt.args.outPath, tt.args.compiler)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRelativePathAdd() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRelativePathAdd() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_GetToolchainFolderPath(t *testing.T) {
	t.Parallel()

	type args struct {
		outPath  string
		compiler string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"./STM32CubeMX", "AC6"}, "STM32CubeMX/MDK-ARM", false},
		{"test2", args{"./STM32CubeMX", "GCC"}, "STM32CubeMX/STM32CubeIDE", false},
		{"test3", args{"./STM32CubeMX", "IAR"}, "STM32CubeMX/EWARM", false},
		{"test4", args{"./STM32CubeMX", "CLANG"}, "STM32CubeMX/STM32CubeIDE", false},
		{"test5", args{"./", "AC6"}, "STM32CubeMX/MDK-ARM", false},
		{"test6", args{"./", "GCC"}, "STM32CubeMX/STM32CubeIDE", false},
		{"test7", args{"./", "IAR"}, "STM32CubeMX/EWARM", false},
		{"test8", args{"./", "CLANG"}, "STM32CubeMX/STM32CubeIDE", false},
		{"fail", args{"", "unknown"}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetToolchainFolderPath(tt.args.outPath, tt.args.compiler)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToolchainFolderPath() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetToolchainFolderPath() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_GetStartupFile(t *testing.T) {
	t.Parallel()

	// Single core
	outPathSC := "../../testdata/testExamples/STM32H7_SC/STM32CubeMX/device"
	var subsystemScAC6 cbuild.SubsystemType
	subsystemScAC6.Compiler = "AC6"
	subsystemScAC6.SubsystemIdx.ProjectType = "single-core"
	subsystemScAC6.SubsystemIdx.ForProjectPart = ""
	var subsystemScGCC cbuild.SubsystemType
	subsystemScGCC.Compiler = "GCC"
	subsystemScGCC.SubsystemIdx.ProjectType = "single-core"
	subsystemScGCC.SubsystemIdx.ForProjectPart = ""
	var subsystemScCLANG cbuild.SubsystemType
	subsystemScCLANG.Compiler = "CLANG"
	subsystemScCLANG.SubsystemIdx.ProjectType = "single-core"
	subsystemScCLANG.SubsystemIdx.ForProjectPart = ""
	var subsystemScIAR cbuild.SubsystemType
	subsystemScIAR.Compiler = "IAR"
	subsystemScIAR.SubsystemIdx.ProjectType = "single-core"
	subsystemScIAR.SubsystemIdx.ForProjectPart = ""

	// Multi core
	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
	var subsystemDcAC6 cbuild.SubsystemType
	subsystemDcAC6.Compiler = "AC6"
	subsystemDcAC6.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcAC6.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcGCC cbuild.SubsystemType
	subsystemDcGCC.Compiler = "GCC"
	subsystemDcGCC.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcGCC.SubsystemIdx.ForProjectPart = "CM7"
	var subsystemDcCLANG cbuild.SubsystemType
	subsystemDcCLANG.Compiler = "CLANG"
	subsystemDcCLANG.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcCLANG.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcIAR cbuild.SubsystemType
	subsystemDcIAR.Compiler = "IAR"
	subsystemDcIAR.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcIAR.SubsystemIdx.ForProjectPart = "CM7"

	// secure nonsecure
	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
	var subsystemTzAC6 cbuild.SubsystemType
	subsystemTzAC6.Compiler = "AC6"
	subsystemTzAC6.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzAC6.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzGCC cbuild.SubsystemType
	subsystemTzGCC.Compiler = "GCC"
	subsystemTzGCC.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzGCC.SubsystemIdx.ForProjectPart = "non-secure"
	var subsystemTzCLANG cbuild.SubsystemType
	subsystemTzCLANG.Compiler = "CLANG"
	subsystemTzCLANG.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzCLANG.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzIAR cbuild.SubsystemType
	subsystemTzIAR.Compiler = "IAR"
	subsystemTzIAR.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzIAR.SubsystemIdx.ForProjectPart = "non-secure"

	// invalid
	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
	var subsystemInv cbuild.SubsystemType
	subsystemInv.Compiler = "AC6"
	subsystemInv.SubsystemIdx.ProjectType = "single-core"
	subsystemInv.SubsystemIdx.ForProjectPart = ""

	type args struct {
		outPath   string
		subsystem cbuild.SubsystemType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test_sc_ac6", args{outPathSC, subsystemScAC6}, filepath.Clean(outPathSC + "/STM32CubeMX/MDK-ARM/startup_stm32h743xx.s"), false},
		{"test_sc_gcc", args{outPathSC, subsystemScGCC}, filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/Application/Startup/startup_stm32h743agix.s"), false},
		{"test_sc_clang", args{outPathSC, subsystemScCLANG}, filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/Application/Startup/startup_stm32h743agix.s"), false},
		{"test_sc_iar", args{outPathSC, subsystemScIAR}, filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/startup_stm32h743xx.s"), false},

		{"test_dc_ac6", args{outPathDC, subsystemDcAC6}, filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/startup_stm32h745xx_CM4.s"), false},
		{"test_dc_gcc", args{outPathDC, subsystemDcGCC}, filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/Application/Startup/startup_stm32h745bgtx.s"), false},
		{"test_dc_clang", args{outPathDC, subsystemDcCLANG}, filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/Application/Startup/startup_stm32h745bgtx.s"), false},
		{"test_dc_iar", args{outPathDC, subsystemDcIAR}, filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/startup_stm32h745xx_CM7.s"), false},

		{"test_tz_ac6", args{outPathTZ, subsystemTzAC6}, filepath.Clean(outPathTZ + "/STM32CubeMX/MDK-ARM/startup_stm32u585xx.s"), false},
		{"test_tz_gcc", args{outPathTZ, subsystemTzGCC}, filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/Application/Startup/startup_stm32u585aiixq.s"), false},
		{"test_tz_clang", args{outPathTZ, subsystemTzCLANG}, filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/Application/Startup/startup_stm32u585aiixq.s"), false},
		{"test_tz_iar", args{outPathTZ, subsystemTzIAR}, filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/startup_stm32u585xx.s"), false},

		{"fail", args{outPathInv, subsystemInv}, "", true}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetStartupFile(tt.args.outPath, &tt.args.subsystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStartupFile() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetStartupFile() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_GetSystemFile(t *testing.T) {
	t.Parallel()

	// Single core
	outPathSC := "../../testdata/testExamples/STM32H7_SC/STM32CubeMX/device"
	var subsystemScAC6 cbuild.SubsystemType
	subsystemScAC6.Compiler = "AC6"
	subsystemScAC6.SubsystemIdx.ProjectType = "single-core"
	subsystemScAC6.SubsystemIdx.ForProjectPart = ""
	var subsystemScGCC cbuild.SubsystemType
	subsystemScGCC.Compiler = "GCC"
	subsystemScGCC.SubsystemIdx.ProjectType = "single-core"
	subsystemScGCC.SubsystemIdx.ForProjectPart = ""
	var subsystemScCLANG cbuild.SubsystemType
	subsystemScCLANG.Compiler = "CLANG"
	subsystemScCLANG.SubsystemIdx.ProjectType = "single-core"
	subsystemScCLANG.SubsystemIdx.ForProjectPart = ""
	var subsystemScIAR cbuild.SubsystemType
	subsystemScIAR.Compiler = "IAR"
	subsystemScIAR.SubsystemIdx.ProjectType = "single-core"
	subsystemScIAR.SubsystemIdx.ForProjectPart = ""

	// Multi core
	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
	var subsystemDcAC6 cbuild.SubsystemType
	subsystemDcAC6.Compiler = "AC6"
	subsystemDcAC6.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcAC6.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcGCC cbuild.SubsystemType
	subsystemDcGCC.Compiler = "GCC"
	subsystemDcGCC.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcGCC.SubsystemIdx.ForProjectPart = "CM7"
	var subsystemDcCLANG cbuild.SubsystemType
	subsystemDcCLANG.Compiler = "CLANG"
	subsystemDcCLANG.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcCLANG.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcIAR cbuild.SubsystemType
	subsystemDcIAR.Compiler = "IAR"
	subsystemDcIAR.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcIAR.SubsystemIdx.ForProjectPart = "CM7"

	// secure nonsecure
	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
	var subsystemTzAC6 cbuild.SubsystemType
	subsystemTzAC6.Compiler = "AC6"
	subsystemTzAC6.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzAC6.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzGCC cbuild.SubsystemType
	subsystemTzGCC.Compiler = "GCC"
	subsystemTzGCC.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzGCC.SubsystemIdx.ForProjectPart = "non-secure"
	var subsystemTzCLANG cbuild.SubsystemType
	subsystemTzCLANG.Compiler = "CLANG"
	subsystemTzCLANG.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzCLANG.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzIAR cbuild.SubsystemType
	subsystemTzIAR.Compiler = "IAR"
	subsystemTzIAR.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzIAR.SubsystemIdx.ForProjectPart = "non-secure"

	// invalid
	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
	var subsystemInv cbuild.SubsystemType
	subsystemInv.Compiler = "AC6"
	subsystemInv.SubsystemIdx.ProjectType = "single-core"
	subsystemInv.SubsystemIdx.ForProjectPart = ""

	type args struct {
		outPath   string
		subsystem cbuild.SubsystemType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test_sc_ac6", args{outPathSC, subsystemScAC6}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_gcc", args{outPathSC, subsystemScGCC}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_clang", args{outPathSC, subsystemScCLANG}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_iar", args{outPathSC, subsystemScIAR}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},

		{"test_dc_ac6", args{outPathDC, subsystemDcAC6}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_gcc", args{outPathDC, subsystemDcGCC}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_clang", args{outPathDC, subsystemDcCLANG}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_iar", args{outPathDC, subsystemDcIAR}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},

		{"test_tz_ac6", args{outPathTZ, subsystemTzAC6}, filepath.Clean(outPathTZ + "/STM32CubeMX/Secure/Src/system_stm32u5xx_s.c"), false},
		{"test_tz_gcc", args{outPathTZ, subsystemTzGCC}, filepath.Clean(outPathTZ + "/STM32CubeMX/NonSecure/Src/system_stm32u5xx_ns.c"), false},
		{"test_tz_clang", args{outPathTZ, subsystemTzCLANG}, filepath.Clean(outPathTZ + "/STM32CubeMX/Secure/Src/system_stm32u5xx_s.c"), false},
		{"test_tz_iar", args{outPathTZ, subsystemTzIAR}, filepath.Clean(outPathTZ + "/STM32CubeMX/NonSecure/Src/system_stm32u5xx_ns.c"), false},

		{"fail", args{outPathInv, subsystemInv}, "", true}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetSystemFile(tt.args.outPath, &tt.args.subsystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemFile() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSystemFile() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_GetLinkerScripts(t *testing.T) {
	t.Parallel()

	// Single core
	outPathSC := "../../testdata/testExamples/STM32H7_SC/STM32CubeMX/device"
	var subsystemScAC6 cbuild.SubsystemType
	subsystemScAC6.Compiler = "AC6"
	subsystemScAC6.SubsystemIdx.ProjectType = "single-core"
	subsystemScAC6.SubsystemIdx.ForProjectPart = ""
	var subsystemScGCC cbuild.SubsystemType
	subsystemScGCC.Compiler = "GCC"
	subsystemScGCC.SubsystemIdx.ProjectType = "single-core"
	subsystemScGCC.SubsystemIdx.ForProjectPart = ""
	var subsystemScCLANG cbuild.SubsystemType
	subsystemScCLANG.Compiler = "CLANG"
	subsystemScCLANG.SubsystemIdx.ProjectType = "single-core"
	subsystemScCLANG.SubsystemIdx.ForProjectPart = ""
	var subsystemScIAR cbuild.SubsystemType
	subsystemScIAR.Compiler = "IAR"
	subsystemScIAR.SubsystemIdx.ProjectType = "single-core"
	subsystemScIAR.SubsystemIdx.ForProjectPart = ""

	// Multi core
	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
	var subsystemDcAC6 cbuild.SubsystemType
	subsystemDcAC6.Compiler = "AC6"
	subsystemDcAC6.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcAC6.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcGCC cbuild.SubsystemType
	subsystemDcGCC.Compiler = "GCC"
	subsystemDcGCC.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcGCC.SubsystemIdx.ForProjectPart = "CM7"
	var subsystemDcCLANG cbuild.SubsystemType
	subsystemDcCLANG.Compiler = "CLANG"
	subsystemDcCLANG.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcCLANG.SubsystemIdx.ForProjectPart = "CM4"
	var subsystemDcIAR cbuild.SubsystemType
	subsystemDcIAR.Compiler = "IAR"
	subsystemDcIAR.SubsystemIdx.ProjectType = "multi-core"
	subsystemDcIAR.SubsystemIdx.ForProjectPart = "CM7"

	// secure nonsecure
	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
	var subsystemTzAC6 cbuild.SubsystemType
	subsystemTzAC6.Compiler = "AC6"
	subsystemTzAC6.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzAC6.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzGCC cbuild.SubsystemType
	subsystemTzGCC.Compiler = "GCC"
	subsystemTzGCC.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzGCC.SubsystemIdx.ForProjectPart = "non-secure"
	var subsystemTzCLANG cbuild.SubsystemType
	subsystemTzCLANG.Compiler = "CLANG"
	subsystemTzCLANG.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzCLANG.SubsystemIdx.ForProjectPart = "secure"
	var subsystemTzIAR cbuild.SubsystemType
	subsystemTzIAR.Compiler = "IAR"
	subsystemTzIAR.SubsystemIdx.ProjectType = "trustzone"
	subsystemTzIAR.SubsystemIdx.ForProjectPart = "non-secure"

	// invalid
	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
	var subsystemInv cbuild.SubsystemType
	subsystemInv.Compiler = "AC6"
	subsystemInv.SubsystemIdx.ProjectType = "single-core"
	subsystemInv.SubsystemIdx.ForProjectPart = ""

	type args struct {
		outPath   string
		subsystem cbuild.SubsystemType
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"test_sc_ac6", args{outPathSC, subsystemScAC6}, nil, false},
		{
			"test_sc_gcc", args{outPathSC, subsystemScGCC},
			[]string{
				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_FLASH.ld"),
				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_RAM.ld"),
			},
			false,
		},
		{
			"test_sc_clang", args{outPathSC, subsystemScCLANG},
			[]string{
				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_FLASH.ld"),
				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_RAM.ld"),
			},
			false,
		},
		{
			"test_sc_iar", args{outPathSC, subsystemScIAR},
			[]string{
				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash.icf"),
				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash_rw_sram1.icf"),
				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash_rw_sram2.icf"),
				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xx_dtcmram.icf"),
				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xx_sram1.icf"),
			},
			false,
		},
		{
			"test_dc_ac6", args{outPathDC, subsystemDcAC6},
			[]string{
				filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/stm32h745xg_flash_CM4.sct"),
				filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/stm32h745xx_sram2_CM4.sct"),
			},
			false,
		},
		{
			"test_dc_gcc", args{outPathDC, subsystemDcGCC},
			[]string{
				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/STM32H745BGTX_FLASH.ld"),
				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/STM32H745BGTX_RAM.ld"),
			},
			false,
		},
		{
			"test_dc_clang", args{outPathDC, subsystemDcCLANG},
			[]string{
				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/STM32H745BGTX_FLASH.ld"),
				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/STM32H745BGTX_RAM.ld"),
			},
			false,
		},
		{
			"test_dc_iar", args{outPathDC, subsystemDcIAR},
			[]string{
				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xg_flash_CM7.icf"),
				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_dtcmram_CM7.icf"),
				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_flash_rw_sram1_CM7.icf"),
				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_sram1_CM7.icf"),
			},
			false,
		},
		{"test_tz_ac6", args{outPathTZ, subsystemTzAC6}, nil, false},
		{
			"test_tz_gcc", args{outPathTZ, subsystemTzGCC},
			[]string{
				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/STM32U585AIIXQ_FLASH.ld"),
				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/STM32U585AIIXQ_RAM.ld"),
			},
			false,
		},
		{
			"test_tz_clang", args{outPathTZ, subsystemTzCLANG},
			[]string{
				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/STM32U585AIIXQ_FLASH.ld"),
				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/STM32U585AIIXQ_RAM.ld"),
			},
			false,
		},
		{
			"test_tz_iar", args{outPathTZ, subsystemTzIAR},
			[]string{
				filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/stm32u585xx_flash_ns.icf"),
				filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/stm32u585xx_sram_ns.icf"),
			},
			false,
		},

		{"fail", args{outPathInv, subsystemInv}, nil, true}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetLinkerScripts(tt.args.outPath, &tt.args.subsystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinkerScripts() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkerScripts() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
