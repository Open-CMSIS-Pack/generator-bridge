/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"path/filepath"
	"testing"

	"github.com/open-cmsis-pack/generator-bridge/internal/cbuild"
)

func Test_GetBridgeInfo(t *testing.T) {

	// Single core Device
	var paramsSC cbuild.ParamsType
	var cgParamsSCTmp cbuild.CbuildGensType
	paramsSC.Board = "BVendorX::BoardY:RevZ"
	paramsSC.Device = "DVendorX::DeviceY"
	paramsSC.ProjectType = "single-core"
	cgParamsSCTmp.Project = "TestProject"
	paramsSC.CbuildGens = append(paramsSC.CbuildGens, cgParamsSCTmp)

	var bParamsSCTmp BridgeParamType
	var bParamsSC []BridgeParamType
	bParamsSCTmp.BoardName = "BoardY"
	bParamsSCTmp.BoardVendor = "BVendorX"
	bParamsSCTmp.Device = "DVendorX::DeviceY"
	bParamsSCTmp.ProjectName = "TestProject"
	bParamsSCTmp.ProjectType = "single-core"
	bParamsSC = append(bParamsSC, bParamsSCTmp)

	// Multi core Device
	var paramsDC cbuild.ParamsType
	var cgParamsDCTmp cbuild.CbuildGensType
	paramsDC.Board = "BoardY"
	paramsDC.Device = "DVendorX::DeviceY"
	paramsDC.ProjectType = "multi-core"
	cgParamsDCTmp.Project = "TestProject1"
	cgParamsDCTmp.ForProjectPart = "CM0P"
	cgParamsDCTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M0+"
	paramsDC.CbuildGens = append(paramsDC.CbuildGens, cgParamsDCTmp)
	cgParamsDCTmp.Project = "TestProject2"
	cgParamsDCTmp.ForProjectPart = "CM4"
	cgParamsDCTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M4"
	paramsDC.CbuildGens = append(paramsDC.CbuildGens, cgParamsDCTmp)

	var bParamsDCTmp BridgeParamType
	var bParamsDC []BridgeParamType
	bParamsDCTmp.BoardName = "BoardY"
	bParamsDCTmp.BoardVendor = ""
	bParamsDCTmp.Device = "DVendorX::DeviceY"
	bParamsDCTmp.ProjectName = "TestProject1"
	bParamsDCTmp.ProjectType = "multi-core"
	bParamsDCTmp.ForProjectPart = "CM0P"
	bParamsDCTmp.CubeContext = "CortexM0Plus"
	bParamsDCTmp.CubeContextFolder = "CM0PLUS"
	bParamsDC = append(bParamsDC, bParamsDCTmp)
	bParamsDCTmp.ProjectName = "TestProject2"
	bParamsDCTmp.ProjectType = "multi-core"
	bParamsDCTmp.ForProjectPart = "CM4"
	bParamsDCTmp.CubeContext = "CortexM4"
	bParamsDCTmp.CubeContextFolder = "CM4"
	bParamsDC = append(bParamsDC, bParamsDCTmp)

	// TZ enabled: Secure Non-Secure
	var paramsTZ cbuild.ParamsType
	var cgParamsTZTmp cbuild.CbuildGensType
	paramsTZ.Board = "BoardY:RevZ"
	paramsTZ.Device = "DeviceY"
	paramsTZ.ProjectType = "trustzone"
	cgParamsTZTmp.Project = "TestProject1"
	cgParamsTZTmp.ForProjectPart = "non-secure"
	cgParamsTZTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M33"
	paramsTZ.CbuildGens = append(paramsTZ.CbuildGens, cgParamsTZTmp)
	cgParamsTZTmp.Project = "TestProject2"
	cgParamsTZTmp.ForProjectPart = "secure"
	cgParamsTZTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M33"
	paramsTZ.CbuildGens = append(paramsTZ.CbuildGens, cgParamsTZTmp)

	var bParamsTZTmp BridgeParamType
	var bParamsTZ []BridgeParamType
	bParamsTZTmp.BoardName = "BoardY"
	bParamsTZTmp.BoardVendor = ""
	bParamsTZTmp.Device = "DeviceY"
	bParamsTZTmp.ProjectName = "TestProject1"
	bParamsTZTmp.ProjectType = "trustzone"
	bParamsTZTmp.ForProjectPart = "non-secure"
	bParamsTZTmp.PairedSecurePart = "TestProject2"
	bParamsTZTmp.CubeContext = "CortexM33NS"
	bParamsTZTmp.CubeContextFolder = "NonSecure"
	bParamsTZ = append(bParamsTZ, bParamsTZTmp)
	bParamsTZTmp.ProjectName = "TestProject2"
	bParamsTZTmp.ForProjectPart = "secure"
	bParamsTZTmp.PairedSecurePart = ""
	bParamsTZTmp.CubeContext = "CortexM33S"
	bParamsTZTmp.CubeContextFolder = "Secure"
	bParamsTZ = append(bParamsTZ, bParamsTZTmp)

	// Boot / Appli
	var paramsBA cbuild.ParamsType
	var cgParamsBATmp cbuild.CbuildGensType
	paramsBA.Board = "DVendorX::BoardY"
	paramsBA.Device = "DVendorX::DeviceY"
	paramsBA.ProjectType = "single-core"
	cgParamsBATmp.Project = "TestProject1"
	cgParamsBATmp.Map = "Appli"
	cgParamsBATmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M7"
	paramsBA.CbuildGens = append(paramsBA.CbuildGens, cgParamsBATmp)
	cgParamsBATmp.Project = "TestProject2"
	cgParamsBATmp.Map = "Boot"
	paramsBA.CbuildGens = append(paramsBA.CbuildGens, cgParamsBATmp)

	var bParamsBATmp BridgeParamType
	var bParamsBA []BridgeParamType
	bParamsBATmp.BoardName = "BoardY"
	bParamsBATmp.BoardVendor = "DVendorX"
	bParamsBATmp.Device = "DVendorX::DeviceY"
	bParamsBATmp.ProjectName = "TestProject1"
	bParamsBATmp.ProjectType = "single-core"
	bParamsBATmp.GeneratorMap = "Appli"
	bParamsBATmp.CubeContext = "Appli"
	bParamsBATmp.CubeContextFolder = "Appli"
	bParamsBA = append(bParamsBA, bParamsBATmp)
	bParamsBATmp.ProjectName = "TestProject2"
	bParamsBATmp.GeneratorMap = "Boot"
	bParamsBATmp.CubeContext = "Boot"
	bParamsBATmp.CubeContextFolder = "Boot"
	bParamsBA = append(bParamsBA, bParamsBATmp)

	// Boot / Appli + Trust Zone
	var paramsBATZ cbuild.ParamsType
	var cgParamsBATZTmp cbuild.CbuildGensType
	paramsBATZ.Board = "DVendorX::BoardY"
	paramsBATZ.Device = "DVendorX::DeviceY"
	paramsBATZ.ProjectType = "trustzone"
	cgParamsBATZTmp.Project = "TestProject1"
	cgParamsBATZTmp.ForProjectPart = "non-secure"
	cgParamsBATZTmp.Map = "AppliNonSecure"
	cgParamsBATZTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M55"
	paramsBATZ.CbuildGens = append(paramsBATZ.CbuildGens, cgParamsBATZTmp)
	cgParamsBATZTmp.Project = "TestProject2"
	cgParamsBATZTmp.ForProjectPart = "secure"
	cgParamsBATZTmp.Map = "AppliSecure"
	cgParamsBATZTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M55"
	paramsBATZ.CbuildGens = append(paramsBATZ.CbuildGens, cgParamsBATZTmp)
	cgParamsBATZTmp.Project = "TestProject3"
	cgParamsBATZTmp.ForProjectPart = "secure"
	cgParamsBATZTmp.Map = "FSBL"
	cgParamsBATZTmp.CbuildGen.BuildGen.Processor.Core = "Cortex-M55"
	paramsBATZ.CbuildGens = append(paramsBATZ.CbuildGens, cgParamsBATZTmp)

	var bParamsBATZTmp BridgeParamType
	var bParamsBATZ []BridgeParamType
	bParamsBATZTmp.BoardName = "BoardY"
	bParamsBATZTmp.BoardVendor = "DVendorX"
	bParamsBATZTmp.Device = "DVendorX::DeviceY"
	bParamsBATZTmp.ProjectName = "TestProject1"
	bParamsBATZTmp.ProjectType = "trustzone"
	bParamsBATZTmp.ForProjectPart = "non-secure"
	bParamsBATZTmp.PairedSecurePart = "TestProject2"
	bParamsBATZTmp.GeneratorMap = "AppliNonSecure"
	bParamsBATZTmp.CubeContext = "AppliNonSecure"
	bParamsBATZTmp.CubeContextFolder = "AppliNonSecure"
	bParamsBATZ = append(bParamsBATZ, bParamsBATZTmp)
	bParamsBATZTmp.ProjectName = "TestProject2"
	bParamsBATZTmp.ForProjectPart = "secure"
	bParamsBATZTmp.PairedSecurePart = ""
	bParamsBATZTmp.GeneratorMap = "AppliSecure"
	bParamsBATZTmp.CubeContext = "AppliSecure"
	bParamsBATZTmp.CubeContextFolder = "AppliSecure"
	bParamsBATZ = append(bParamsBATZ, bParamsBATZTmp)
	bParamsBATZTmp.ProjectName = "TestProject3"
	bParamsBATZTmp.ForProjectPart = "secure"
	bParamsBATZTmp.PairedSecurePart = ""
	bParamsBATZTmp.GeneratorMap = "FSBL"
	bParamsBATZTmp.CubeContext = "FSBL"
	bParamsBATZTmp.CubeContextFolder = "FSBL"
	bParamsBATZ = append(bParamsBATZ, bParamsBATZTmp)

	type args struct {
		params *cbuild.ParamsType
	}
	tests := []struct {
		name    string
		args    args
		want    []BridgeParamType
		wantErr error
	}{
		{"testSC", args{&paramsSC}, bParamsSC, nil},
		{"testDC", args{&paramsDC}, bParamsDC, nil},
		{"testTZ", args{&paramsTZ}, bParamsTZ, nil},
		{"testBA", args{&paramsBA}, bParamsBA, nil},
		{"testBATZ", args{&paramsBATZ}, bParamsBATZ, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var retBridgeParams []BridgeParamType
			err := GetBridgeInfo(tt.args.params, &retBridgeParams)
			if err != tt.wantErr {
				t.Errorf("GetBridgeInfo() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if len(retBridgeParams) != len(tt.want) {
				t.Errorf("GetBridgeInfo() %s: Un-expected length of returned BridgeParams", tt.name)
			} else {
				for i := range tt.want {
					if retBridgeParams[i].BoardName != tt.want[i].BoardName {
						t.Errorf("GetBridgeInfo() %s BoardName = %v, want %v", tt.name, retBridgeParams[i].BoardName, tt.want[i].BoardName)
					}
					if retBridgeParams[i].BoardVendor != tt.want[i].BoardVendor {
						t.Errorf("GetBridgeInfo() %s BoardVendor = %v, want %v", tt.name, retBridgeParams[i].BoardVendor, tt.want[i].BoardVendor)
					}
					if retBridgeParams[i].Device != tt.want[i].Device {
						t.Errorf("GetBridgeInfo() %s Device = %v, want %v", tt.name, retBridgeParams[i].Device, tt.want[i].Device)
					}
					if retBridgeParams[i].Output != tt.want[i].Output {
						t.Errorf("GetBridgeInfo() %s Output = %v, want %v", tt.name, retBridgeParams[i].Output, tt.want[i].Output)
					}
					if retBridgeParams[i].ProjectName != tt.want[i].ProjectName {
						t.Errorf("GetBridgeInfo() %s ProjectName = %v, want %v", tt.name, retBridgeParams[i].ProjectName, tt.want[i].ProjectName)
					}
					if retBridgeParams[i].ProjectType != tt.want[i].ProjectType {
						t.Errorf("GetBridgeInfo() %s ProjectType = %v, want %v", tt.name, retBridgeParams[i].ProjectType, tt.want[i].ProjectType)
					}
					if retBridgeParams[i].ForProjectPart != tt.want[i].ForProjectPart {
						t.Errorf("GetBridgeInfo() %s ForProjectPart = %v, want %v", tt.name, retBridgeParams[i].ForProjectPart, tt.want[i].ForProjectPart)
					}
					if retBridgeParams[i].PairedSecurePart != tt.want[i].PairedSecurePart {
						t.Errorf("GetBridgeInfo() %s PairedSecurePart = %v, want %v", tt.name, retBridgeParams[i].PairedSecurePart, tt.want[i].PairedSecurePart)
					}
					if retBridgeParams[i].Compiler != tt.want[i].Compiler {
						t.Errorf("GetBridgeInfo() %s Compiler = %v, want %v", tt.name, retBridgeParams[i].Compiler, tt.want[i].Compiler)
					}
					if retBridgeParams[i].GeneratorMap != tt.want[i].GeneratorMap {
						t.Errorf("GetBridgeInfo() %s GeneratorMap = %v, want %v", tt.name, retBridgeParams[i].GeneratorMap, tt.want[i].GeneratorMap)
					}
					if retBridgeParams[i].CgenName != tt.want[i].CgenName {
						t.Errorf("GetBridgeInfo() %s CgenName = %v, want %v", tt.name, retBridgeParams[i].CgenName, tt.want[i].CgenName)
					}
					if retBridgeParams[i].CubeContext != tt.want[i].CubeContext {
						t.Errorf("GetBridgeInfo() %s CubeContext = %v, want %v", tt.name, retBridgeParams[i].CubeContext, tt.want[i].CubeContext)
					}
					if retBridgeParams[i].CubeContextFolder != tt.want[i].CubeContextFolder {
						t.Errorf("GetBridgeInfo() %s CubeContextFolder = %v, want %v", tt.name, retBridgeParams[i].CubeContextFolder, tt.want[i].CubeContextFolder)
					}
				}
			}
		})
	}
}

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
			got = filepath.ToSlash(got)
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
			got = filepath.ToSlash(got)
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
	var infoScAC6 BridgeParamType
	infoScAC6.Compiler = "AC6"
	infoScAC6.ProjectType = "single-core"
	infoScAC6.ForProjectPart = ""
	infoScAC6.CubeContext = ""
	infoScAC6.CubeContextFolder = ""
	var infoScGCC BridgeParamType
	infoScGCC.Compiler = "GCC"
	infoScGCC.ProjectType = "single-core"
	infoScGCC.ForProjectPart = ""
	infoScGCC.CubeContext = ""
	infoScGCC.CubeContextFolder = ""
	var infoScCLANG BridgeParamType
	infoScCLANG.Compiler = "CLANG"
	infoScCLANG.ProjectType = "single-core"
	infoScCLANG.ForProjectPart = ""
	infoScCLANG.CubeContext = ""
	infoScCLANG.CubeContextFolder = ""
	var infoScIAR BridgeParamType
	infoScIAR.Compiler = "IAR"
	infoScIAR.ProjectType = "single-core"
	infoScIAR.ForProjectPart = ""
	infoScIAR.CubeContext = ""
	infoScIAR.CubeContextFolder = ""

	// Multi core
	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
	var infoDcAC6 BridgeParamType
	infoDcAC6.Compiler = "AC6"
	infoDcAC6.ProjectType = "multi-core"
	infoDcAC6.ForProjectPart = "CM4"
	infoDcAC6.CubeContext = "CortexM4"
	infoDcAC6.CubeContextFolder = "CM4"
	var infoDcGCC BridgeParamType
	infoDcGCC.Compiler = "GCC"
	infoDcGCC.ProjectType = "multi-core"
	infoDcGCC.ForProjectPart = "CM7"
	infoDcGCC.CubeContext = "CortexM7"
	infoDcGCC.CubeContextFolder = "CM7"
	var infoDcCLANG BridgeParamType
	infoDcCLANG.Compiler = "CLANG"
	infoDcCLANG.ProjectType = "multi-core"
	infoDcCLANG.ForProjectPart = "CM4"
	infoDcCLANG.CubeContext = "CortexM4"
	infoDcCLANG.CubeContextFolder = "CM4"
	var infoDcIAR BridgeParamType
	infoDcIAR.Compiler = "IAR"
	infoDcIAR.ProjectType = "multi-core"
	infoDcIAR.ForProjectPart = "CM7"
	infoDcIAR.CubeContext = "CortexM7"
	infoDcIAR.CubeContextFolder = "CM7"

	// secure nonsecure
	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
	var infoTzAC6 BridgeParamType
	infoTzAC6.Compiler = "AC6"
	infoTzAC6.ProjectType = "trustzone"
	infoTzAC6.ForProjectPart = "secure"
	infoTzAC6.CubeContext = "CortexM33S"
	infoTzAC6.CubeContextFolder = "Secure"
	var infoTzGCC BridgeParamType
	infoTzGCC.Compiler = "GCC"
	infoTzGCC.ProjectType = "trustzone"
	infoTzGCC.ForProjectPart = "non-secure"
	infoTzGCC.CubeContext = "CortexM33NS"
	infoTzGCC.CubeContextFolder = "NonSecure"
	var infoTzCLANG BridgeParamType
	infoTzCLANG.Compiler = "CLANG"
	infoTzCLANG.ProjectType = "trustzone"
	infoTzCLANG.ForProjectPart = "secure"
	infoTzCLANG.CubeContext = "CortexM33S"
	infoTzCLANG.CubeContextFolder = "Secure"
	var infoTzIAR BridgeParamType
	infoTzIAR.Compiler = "IAR"
	infoTzIAR.ProjectType = "trustzone"
	infoTzIAR.ForProjectPart = "non-secure"
	infoTzIAR.CubeContext = "CortexM33NS"
	infoTzIAR.CubeContextFolder = "NonSecure"

	// Multi core (M4 & M0+)
	outPathDCM0P := "../../testdata/testExamples/STM32WL_DC/test/STM32CubeMX/STM32WL54CCUx"
	var infoDcM0PAC6 BridgeParamType
	infoDcM0PAC6.Compiler = "AC6"
	infoDcM0PAC6.ProjectType = "multi-core"
	infoDcM0PAC6.ForProjectPart = "CM0P"
	infoDcM0PAC6.CubeContext = "CortexM0Plus"
	infoDcM0PAC6.CubeContextFolder = "CM0PLUS"
	var infoDcM0PGCC BridgeParamType
	infoDcM0PGCC.Compiler = "GCC"
	infoDcM0PGCC.ProjectType = "multi-core"
	infoDcM0PGCC.ForProjectPart = "CM4"
	infoDcM0PGCC.CubeContext = "CortexM4"
	infoDcM0PGCC.CubeContextFolder = "CM4"
	var infoDcM0PCLANG BridgeParamType
	infoDcM0PCLANG.Compiler = "CLANG"
	infoDcM0PCLANG.ProjectType = "multi-core"
	infoDcM0PCLANG.ForProjectPart = "CM0P"
	infoDcM0PCLANG.CubeContext = "CortexM0Plus"
	infoDcM0PCLANG.CubeContextFolder = "CM0PLUS"
	var infoDcM0PIAR BridgeParamType
	infoDcM0PIAR.Compiler = "IAR"
	infoDcM0PIAR.ProjectType = "multi-core"
	infoDcM0PIAR.ForProjectPart = "CM4"
	infoDcM0PIAR.CubeContext = "CortexM4"
	infoDcM0PIAR.CubeContextFolder = "CM4"

	// invalid
	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
	var infoInv BridgeParamType
	infoInv.Compiler = "AC6"
	infoInv.ProjectType = "single-core"
	infoInv.ForProjectPart = ""

	type args struct {
		outPath string
		info    BridgeParamType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test_sc_ac6", args{outPathSC, infoScAC6}, filepath.Clean(outPathSC + "/STM32CubeMX/MDK-ARM/startup_stm32h743xx.s"), false},
		{"test_sc_gcc", args{outPathSC, infoScGCC}, filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/Application/Startup/startup_stm32h743agix.s"), false},
		{"test_sc_clang", args{outPathSC, infoScCLANG}, filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/Application/Startup/startup_stm32h743agix.s"), false},
		{"test_sc_iar", args{outPathSC, infoScIAR}, filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/startup_stm32h743xx.s"), false},

		{"test_dc_ac6", args{outPathDC, infoDcAC6}, filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/startup_stm32h745xx_CM4.s"), false},
		{"test_dc_gcc", args{outPathDC, infoDcGCC}, filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/Application/Startup/startup_stm32h745bgtx.s"), false},
		{"test_dc_clang", args{outPathDC, infoDcCLANG}, filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/Application/Startup/startup_stm32h745bgtx.s"), false},
		{"test_dc_iar", args{outPathDC, infoDcIAR}, filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/startup_stm32h745xx_CM7.s"), false},

		{"test_tz_ac6", args{outPathTZ, infoTzAC6}, filepath.Clean(outPathTZ + "/STM32CubeMX/MDK-ARM/startup_stm32u585xx.s"), false},
		{"test_tz_gcc", args{outPathTZ, infoTzGCC}, filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/Application/Startup/startup_stm32u585aiixq.s"), false},
		{"test_tz_clang", args{outPathTZ, infoTzCLANG}, filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/Application/Startup/startup_stm32u585aiixq.s"), false},
		{"test_tz_iar", args{outPathTZ, infoTzIAR}, filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/startup_stm32u585xx.s"), false},

		{"test_tz_ac6", args{outPathDCM0P, infoDcM0PAC6}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/MDK-ARM/startup_stm32wl54xx_cm0plus.s"), false},
		{"test_tz_gcc", args{outPathDCM0P, infoDcM0PGCC}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/STM32CubeIDE/CM4/Application/Startup/startup_stm32wl54ccux.s"), false},
		{"test_tz_clang", args{outPathDCM0P, infoDcM0PCLANG}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/STM32CubeIDE/CM0PLUS/Application/Startup/startup_stm32wl54ccux.s"), false},
		{"test_tz_iar", args{outPathDCM0P, infoDcM0PIAR}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/EWARM/startup_stm32wl54xx_cm4.s"), false},

		{"fail", args{outPathInv, infoInv}, "", true}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetStartupFile(tt.args.outPath, tt.args.info)
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
	var infoScAC6 BridgeParamType
	infoScAC6.Compiler = "AC6"
	infoScAC6.ProjectType = "single-core"
	infoScAC6.ForProjectPart = ""
	infoScAC6.CubeContext = ""
	infoScAC6.CubeContextFolder = ""
	var infoScGCC BridgeParamType
	infoScGCC.Compiler = "GCC"
	infoScGCC.ProjectType = "single-core"
	infoScGCC.ForProjectPart = ""
	infoScGCC.CubeContext = ""
	infoScGCC.CubeContextFolder = ""
	var infoScCLANG BridgeParamType
	infoScCLANG.Compiler = "CLANG"
	infoScCLANG.ProjectType = "single-core"
	infoScCLANG.ForProjectPart = ""
	infoScCLANG.CubeContext = ""
	infoScCLANG.CubeContextFolder = ""
	var infoScIAR BridgeParamType
	infoScIAR.Compiler = "IAR"
	infoScIAR.ProjectType = "single-core"
	infoScIAR.ForProjectPart = ""
	infoScIAR.CubeContext = ""
	infoScIAR.CubeContextFolder = ""

	// Multi core
	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
	var infoDcAC6 BridgeParamType
	infoDcAC6.Compiler = "AC6"
	infoDcAC6.ProjectType = "multi-core"
	infoDcAC6.ForProjectPart = "CM4"
	infoDcAC6.CubeContext = "CortexM4"
	infoDcAC6.CubeContextFolder = "CM4"
	var infoDcGCC BridgeParamType
	infoDcGCC.Compiler = "GCC"
	infoDcGCC.ProjectType = "multi-core"
	infoDcGCC.ForProjectPart = "CM7"
	infoDcGCC.CubeContext = "CortexM7"
	infoDcGCC.CubeContextFolder = "CM7"
	var infoDcCLANG BridgeParamType
	infoDcCLANG.Compiler = "CLANG"
	infoDcCLANG.ProjectType = "multi-core"
	infoDcCLANG.ForProjectPart = "CortexM4"
	infoDcCLANG.CubeContext = "CM4"
	infoDcCLANG.CubeContextFolder = "CM4"
	var infoDcIAR BridgeParamType
	infoDcIAR.Compiler = "IAR"
	infoDcIAR.ProjectType = "multi-core"
	infoDcIAR.ForProjectPart = "CM7"
	infoDcIAR.CubeContext = "CortexM7"
	infoDcIAR.CubeContextFolder = "CM7"

	// secure nonsecure
	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
	var infoTzAC6 BridgeParamType
	infoTzAC6.Compiler = "AC6"
	infoTzAC6.ProjectType = "trustzone"
	infoTzAC6.ForProjectPart = "secure"
	infoTzAC6.CubeContext = "CortexM33S"
	infoTzAC6.CubeContextFolder = "Secure"
	var infoTzGCC BridgeParamType
	infoTzGCC.Compiler = "GCC"
	infoTzGCC.ProjectType = "trustzone"
	infoTzGCC.ForProjectPart = "non-secure"
	infoTzGCC.CubeContext = "CortexM33NS"
	infoTzGCC.CubeContextFolder = "NonSecure"
	var infoTzCLANG BridgeParamType
	infoTzCLANG.Compiler = "CLANG"
	infoTzCLANG.ProjectType = "trustzone"
	infoTzCLANG.ForProjectPart = "secure"
	infoTzCLANG.CubeContext = "CortexM33S"
	infoTzCLANG.CubeContextFolder = "Secure"
	var infoTzIAR BridgeParamType
	infoTzIAR.Compiler = "IAR"
	infoTzIAR.ProjectType = "trustzone"
	infoTzIAR.ForProjectPart = "non-secure"
	infoTzIAR.CubeContext = "CortexM33NS"
	infoTzIAR.CubeContextFolder = "NonSecure"

	// Multi core (M4 & M0+)
	outPathDCM0P := "../../testdata/testExamples/STM32WL_DC/test/STM32CubeMX/STM32WL54CCUx"
	var infoDcM0PAC6 BridgeParamType
	infoDcM0PAC6.Compiler = "AC6"
	infoDcM0PAC6.ProjectType = "multi-core"
	infoDcM0PAC6.ForProjectPart = "CM0P"
	infoDcM0PAC6.CubeContext = "CortexM0Plus"
	infoDcM0PAC6.CubeContextFolder = "CM0PLUS"
	var infoDcM0PGCC BridgeParamType
	infoDcM0PGCC.Compiler = "GCC"
	infoDcM0PGCC.ProjectType = "multi-core"
	infoDcM0PGCC.ForProjectPart = "CM4"
	infoDcM0PGCC.CubeContext = "CortexM4"
	infoDcM0PGCC.CubeContextFolder = "CM4"
	var infoDcM0PCLANG BridgeParamType
	infoDcM0PCLANG.Compiler = "CLANG"
	infoDcM0PCLANG.ProjectType = "multi-core"
	infoDcM0PCLANG.ForProjectPart = "CM0P"
	infoDcM0PCLANG.CubeContext = "CortexM0Plus"
	infoDcM0PCLANG.CubeContextFolder = "CM0PLUS"
	var infoDcM0PIAR BridgeParamType
	infoDcM0PIAR.Compiler = "IAR"
	infoDcM0PIAR.ProjectType = "multi-core"
	infoDcM0PIAR.ForProjectPart = "CM4"
	infoDcM0PIAR.CubeContext = "CortexM4"
	infoDcM0PIAR.CubeContextFolder = "CM4"

	// invalid
	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
	var infoInv BridgeParamType
	infoInv.Compiler = "AC6"
	infoInv.ProjectType = "single-core"
	infoInv.ForProjectPart = ""

	type args struct {
		outPath string
		info    BridgeParamType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test_sc_ac6", args{outPathSC, infoScAC6}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_gcc", args{outPathSC, infoScGCC}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_clang", args{outPathSC, infoScCLANG}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},
		{"test_sc_iar", args{outPathSC, infoScIAR}, filepath.Clean(outPathSC + "/STM32CubeMX/Src/system_stm32h7xx.c"), false},

		{"test_dc_ac6", args{outPathDC, infoDcAC6}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_gcc", args{outPathDC, infoDcGCC}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_clang", args{outPathDC, infoDcCLANG}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},
		{"test_dc_iar", args{outPathDC, infoDcIAR}, filepath.Clean(outPathDC + "/STM32CubeMX/Common/Src/system_stm32h7xx_dualcore_boot_cm4_cm7.c"), false},

		{"test_tz_ac6", args{outPathTZ, infoTzAC6}, filepath.Clean(outPathTZ + "/STM32CubeMX/Secure/Src/system_stm32u5xx_s.c"), false},
		{"test_tz_gcc", args{outPathTZ, infoTzGCC}, filepath.Clean(outPathTZ + "/STM32CubeMX/NonSecure/Src/system_stm32u5xx_ns.c"), false},
		{"test_tz_clang", args{outPathTZ, infoTzCLANG}, filepath.Clean(outPathTZ + "/STM32CubeMX/Secure/Src/system_stm32u5xx_s.c"), false},
		{"test_tz_iar", args{outPathTZ, infoTzIAR}, filepath.Clean(outPathTZ + "/STM32CubeMX/NonSecure/Src/system_stm32u5xx_ns.c"), false},

		{"test_tz_ac6", args{outPathDCM0P, infoDcM0PAC6}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/Common/System/system_stm32wlxx.c"), false},
		{"test_tz_gcc", args{outPathDCM0P, infoDcM0PGCC}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/Common/System/system_stm32wlxx.c"), false},
		{"test_tz_clang", args{outPathDCM0P, infoDcM0PCLANG}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/Common/System/system_stm32wlxx.c"), false},
		{"test_tz_iar", args{outPathDCM0P, infoDcM0PIAR}, filepath.Clean(outPathDCM0P + "/STM32CubeMX/Common/System/system_stm32wlxx.c"), false},

		{"fail", args{outPathInv, infoInv}, "", true}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetSystemFile(tt.args.outPath, tt.args.info)
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

// func Test_GetLinkerScripts(t *testing.T) {
// 	t.Parallel()

// 	// Single core
// 	outPathSC := "../../testdata/testExamples/STM32H7_SC/STM32CubeMX/device"
// 	var infoScAC6 BridgeParamType
// 	infoScAC6.Compiler = "AC6"
// 	infoScAC6.ProjectType = "single-core"
// 	infoScAC6.ForProjectPart = ""
// 	var infoScGCC BridgeParamType
// 	infoScGCC.Compiler = "GCC"
// 	infoScGCC.ProjectType = "single-core"
// 	infoScGCC.ForProjectPart = ""
// 	var infoScCLANG BridgeParamType
// 	infoScCLANG.Compiler = "CLANG"
// 	infoScCLANG.ProjectType = "single-core"
// 	infoScCLANG.ForProjectPart = ""
// 	var infoScIAR BridgeParamType
// 	infoScIAR.Compiler = "IAR"
// 	infoScIAR.ProjectType = "single-core"
// 	infoScIAR.ForProjectPart = ""

// 	// Multi core
// 	outPathDC := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx"
// 	var infoDcAC6 BridgeParamType
// 	infoDcAC6.Compiler = "AC6"
// 	infoDcAC6.ProjectType = "multi-core"
// 	infoDcAC6.ForProjectPart = "CM4"
// 	var infoDcGCC BridgeParamType
// 	infoDcGCC.Compiler = "GCC"
// 	infoDcGCC.ProjectType = "multi-core"
// 	infoDcGCC.ForProjectPart = "CM7"
// 	var infoDcCLANG BridgeParamType
// 	infoDcCLANG.Compiler = "CLANG"
// 	infoDcCLANG.ProjectType = "multi-core"
// 	infoDcCLANG.ForProjectPart = "CM4"
// 	var infoDcIAR BridgeParamType
// 	infoDcIAR.Compiler = "IAR"
// 	infoDcIAR.ProjectType = "multi-core"
// 	infoDcIAR.ForProjectPart = "CM7"

// 	// secure nonsecure
// 	outPathTZ := "../../testdata/testExamples/STM32U5_TZ/STM32CubeMX/Board"
// 	var infoTzAC6 BridgeParamType
// 	infoTzAC6.Compiler = "AC6"
// 	infoTzAC6.ProjectType = "trustzone"
// 	infoTzAC6.ForProjectPart = "secure"
// 	var infoTzGCC BridgeParamType
// 	infoTzGCC.Compiler = "GCC"
// 	infoTzGCC.ProjectType = "trustzone"
// 	infoTzGCC.ForProjectPart = "non-secure"
// 	var infoTzCLANG BridgeParamType
// 	infoTzCLANG.Compiler = "CLANG"
// 	infoTzCLANG.ProjectType = "trustzone"
// 	infoTzCLANG.ForProjectPart = "secure"
// 	var infoTzIAR BridgeParamType
// 	infoTzIAR.Compiler = "IAR"
// 	infoTzIAR.ProjectType = "trustzone"
// 	infoTzIAR.ForProjectPart = "non-secure"

// 	// invalid
// 	outPathInv := "../../testdata/testExamples/STM32H7_DC/STM32CubeMX/STM32H745BGTx/invalid_folder"
// 	var infoInv BridgeParamType
// 	infoInv.Compiler = "AC6"
// 	infoInv.ProjectType = "single-core"
// 	infoInv.ForProjectPart = ""

// 	type args struct {
// 		outPath string
// 		info    BridgeParamType
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []string
// 		wantErr bool
// 	}{
// 		{"test_sc_ac6", args{outPathSC, infoScAC6}, nil, false},
// 		{
// 			"test_sc_gcc", args{outPathSC, infoScGCC},
// 			[]string{
// 				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_FLASH.ld"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_sc_clang", args{outPathSC, infoScCLANG},
// 			[]string{
// 				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_FLASH.ld"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/STM32CubeIDE/STM32H743AGIX_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_sc_iar", args{outPathSC, infoScIAR},
// 			[]string{
// 				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash.icf"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash_rw_sram1.icf"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xg_flash_rw_sram2.icf"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xx_dtcmram.icf"),
// 				filepath.Clean(outPathSC + "/STM32CubeMX/EWARM/stm32h743xx_sram1.icf"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_dc_ac6", args{outPathDC, infoDcAC6},
// 			[]string{
// 				filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/stm32h745xg_flash_CM4.sct"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/MDK-ARM/stm32h745xx_sram2_CM4.sct"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_dc_gcc", args{outPathDC, infoDcGCC},
// 			[]string{
// 				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/STM32H745BGTX_FLASH.ld"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM7/STM32H745BGTX_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_dc_clang", args{outPathDC, infoDcCLANG},
// 			[]string{
// 				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/STM32H745BGTX_FLASH.ld"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/STM32CubeIDE/CM4/STM32H745BGTX_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_dc_iar", args{outPathDC, infoDcIAR},
// 			[]string{
// 				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xg_flash_CM7.icf"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_dtcmram_CM7.icf"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_flash_rw_sram1_CM7.icf"),
// 				filepath.Clean(outPathDC + "/STM32CubeMX/EWARM/stm32h745xx_sram1_CM7.icf"),
// 			},
// 			false,
// 		},
// 		{"test_tz_ac6", args{outPathTZ, infoTzAC6}, nil, false},
// 		{
// 			"test_tz_gcc", args{outPathTZ, infoTzGCC},
// 			[]string{
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/STM32U585AIIXQ_FLASH.ld"),
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/NonSecure/STM32U585AIIXQ_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_tz_clang", args{outPathTZ, infoTzCLANG},
// 			[]string{
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/STM32U585AIIXQ_FLASH.ld"),
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/STM32CubeIDE/Secure/STM32U585AIIXQ_RAM.ld"),
// 			},
// 			false,
// 		},
// 		{
// 			"test_tz_iar", args{outPathTZ, infoTzIAR},
// 			[]string{
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/stm32u585xx_flash_ns.icf"),
// 				filepath.Clean(outPathTZ + "/STM32CubeMX/EWARM/stm32u585xx_sram_ns.icf"),
// 			},
// 			false,
// 		},

// 		{"fail", args{outPathInv, infoInv}, nil, true}}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got, err := GetLinkerScripts(tt.args.outPath, tt.args.info)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetLinkerScripts() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetLinkerScripts() %s = %v, want %v", tt.name, got, tt.want)
// 			}
// 		})
// 	}
// }
