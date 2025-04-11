/*
 * Copyright (c) 2023-2024 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_MXDeviceContexts(t *testing.T) {

	tests := []struct {
		name       string
		cfgPath    string
		cfgName    string
		cfgRefName string
		wantErr    bool
	}{
		{"STM32H7_SC", "STM32H7_SC/STM32CubeMX/device/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32H7_DC:CM4", "STM32H7_DC/STM32CubeMX/STM32H745BGTx/MX_Device/CM4/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32H7_DC:CM7", "STM32H7_DC/STM32CubeMX/STM32H745BGTx/MX_Device/CM7/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32U5_noTZ", "STM32U5_noTZ/STM32CubeMX/Board/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32U5_TZ:NonSecure", "STM32U5_TZ/STM32CubeMX/Board/MX_Device/NonSecure/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32U5_TZ:Secure", "STM32U5_TZ/STM32CubeMX/Board/MX_Device/Secure/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32WL_DC:CM0PLUS", "STM32WL_DC/test/STM32CubeMX/STM32WL54CCUx/MX_Device/CM0PLUS/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32WL_DC:CM4", "STM32WL_DC/test/STM32CubeMX/STM32WL54CCUx/MX_Device/CM4/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32H5", "STM32H5/STM32CubeMX/STM32H573IIKxQ/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32H7", "STM32H7/STM32CubeMX/STM32H743XIHx/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32F7", "STM32F7/STM32CubeMX/STM32F746NGHx/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32U5", "STM32U5/STM32CubeMX/STM32U5G9ZJTxQ/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32F2", "STM32F2/STM32CubeMX/STM32F217IGHx/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32F4", "STM32F4/STM32CubeMX/STM32F469NIHx/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
		{"STM32G4", "STM32G4/STM32CubeMX/STM32G474QETx/MX_Device/", "MX_Device.h", "MX_Device_ref.h", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.cfgPath = filepath.Join("../../testdata/testExamples/", tt.cfgPath)

			data, err := os.ReadFile(filepath.Join(tt.cfgPath, tt.cfgName))
			if err != nil {
				t.Errorf(" %s; cannot open MX_Device.h file", tt.name)
				return
			}

			data_ref, err_ref := os.ReadFile(filepath.Join(tt.cfgPath, tt.cfgRefName))
			if err_ref != nil {
				t.Errorf(" %s; cannot open MX_Device_ref.h file", tt.name)
				return
			}

			// Normalize line endings inline
			content := strings.ReplaceAll(string(data), "\r\n", "\n")
			content_ref := strings.ReplaceAll(string(data_ref), "\r\n", "\n")

			lines := strings.Split(content, "\n")
			lines_ref := strings.Split(content_ref, "\n")

			content = strings.Join(lines[3:], "\n")
			content_ref = strings.Join(lines_ref[3:], "\n")

			if reflect.DeepEqual((content != content_ref), !tt.wantErr) {
				t.Errorf(" %s; MX_Device.h file content mismatch", tt.name)
			}
		})
	}
}

func Test_ReadContexts(t *testing.T) {

	tests := []struct {
		name    string
		iocFile string
		params  BridgeParamType
		wantErr bool
	}{
		{"STM32H7_SC", "../../testdata/testExamples/STM32H7_SC/STM32CubeMX/device/STM32CubeMX/STM32CubeMX.ioc", BridgeParamType{"device", "", "STM32H743AGIx", "", "test", "single-core", "", "", "AC6", "", "test", "", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argsSlice := []BridgeParamType{tt.params}
			if err := ReadContexts(tt.iocFile, argsSlice); (err != nil) != tt.wantErr {
				t.Errorf("ReadContexts() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func Test_createContextMap(t *testing.T) {
	t.Parallel()

	var parts = make(map[string]map[string]string)
	parts["Mcu"] = map[string]string{"Family": "STM32U5"}
	parts["PA10"] = map[string]string{"GPIOParameters": "GPIO_Label"}

	type args struct {
		iocFile string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]string
		wantErr bool
	}{
		{"test.ioc", args{"../../testdata/stm32cubemx/test.ioc"}, parts, false},
		{"nix", args{"xxx"}, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := createContextMap(tt.args.iocFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("createContextMap() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createContextMap() %s got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_writeMXdeviceH(t *testing.T) {
	var mcuContext0 = make(map[string]map[string]string)
	mcuContext0["Mcu"] = map[string]string{"nixIPs": "myContext"}
	var mcuContext1 = make(map[string]map[string]string)
	mcuContext1["Mcu"] = map[string]string{"IPs": "myContext"}
	var mcuContext2 = make(map[string]map[string]string)
	mcuContext2["Mcu"] = map[string]string{"IP1": "SPI2"}
	mcuContext2["SPI2"] = map[string]string{"VirtualModex": "ModeX"}
	mcuContext2["P1"] = map[string]string{"Signal": "SPI2"}

	type args struct {
		contextMap map[string]map[string]string
		srcFolder  string
		mspName    string
		cfgPath    string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Mcu", args{mcuContext2, "../../testdata/stm32cubemx", "test_msp.c", "../../testdata/stm32cubemx/cfg", ""}, false},
		{"wrong context", args{mcuContext0, "../../testdata/stm32cubemx", "test_msp.c", "../../testdata/stm32cubemx/cfg", "context"}, true},
		{"Mcu Context", args{mcuContext1, "../../testdata/stm32cubemx", "test_msp.c", "../../testdata/stm32cubemx/cfg", "Mcu"}, false},
		{"wrong myContext", args{mcuContext1, "../../testdata/stm32cubemx", "test_msp.c", "../../testdata/stm32cubemx/cfg", "context"}, true},
		{"wrong msp", args{mcuContext1, "", "msp", "../../testdata/stm32cubemx/cfg", "context"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.args.srcFolder + "/../" + tt.args.cfgPath)
			if err := writeMXdeviceH(tt.args.contextMap, tt.args.srcFolder, tt.args.mspName, tt.args.cfgPath, tt.args.context); (err != nil) != tt.wantErr {
				t.Errorf("writeMXdeviceH() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func Test_getContexts(t *testing.T) {
	t.Parallel()

	var mcuContext1 = make(map[string]map[string]string)
	mcuContext1["Mcu"] = map[string]string{"ContextTest": "myContext"}

	var mcuContext2 = make(map[string]map[string]string)
	mcuContext2["Mcu"] = map[string]string{"Context1": "myContext1", "Context2": "myContext2"}

	var oneEmpty = make(map[int]string)
	oneEmpty[0] = ""

	var two1 = make(map[int]string)
	two1[1] = "myContext1"
	two1[2] = "myContext2"

	var two2 = make(map[int]string)
	two2[2] = "myContext1"
	two2[1] = "myContext2"

	type args struct {
		contextMap map[string]map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want1   map[int]string
		want2   map[int]string
		wantErr bool
	}{
		{"1 Context", args{mcuContext1}, oneEmpty, oneEmpty, false},
		{"2 Contexts", args{mcuContext2}, two1, two2, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getContexts(tt.args.contextMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getContexts() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want1) && !reflect.DeepEqual(got, tt.want2) {
				t.Errorf("getContexts() %s = %v, want %v", tt.name, got, tt.want1)
			}
		})
	}
}

func Test_getPeripherals(t *testing.T) {
	//	t.Parallel()

	var parts = make(map[string]map[string]string)
	parts["Mcu"] = map[string]string{"IP1": "UART", "IP2": "xx", "IP3": "USB"}
	parts["Context1"] = map[string]string{"IPs": "CORTEX_M4\\:I,UART5\\:I,USB_DEVICE_M4\\:I"}

	var parts1 = make(map[string]map[string]string)
	parts1["ccc"] = map[string]string{}

	var parts2 = make(map[string]map[string]string)
	parts2["Context1"] = map[string]string{"xxx": "jdsfhkha"}

	var parts3 = make(map[string]map[string]string)
	parts3["Context1"] = map[string]string{"IPs": "CORTEX_M4\\:I,UART5\\:I,USB_DEVICE_M4\\:I"}

	type args struct {
		contextMap map[string]map[string]string
		context    string
	}
	tests := []struct {
		name    string
		args    args
		want1   []string
		want2   []string
		wantErr bool
	}{
		{"test", args{parts, "Context1"}, []string{"UART", "USB"}, []string{"USB", "UART"}, false},
		{"fail1", args{parts1, "xxx"}, nil, nil, true},
		{"fail2", args{parts2, "Context1"}, nil, nil, true},
		{"fail3", args{parts3, "Context1"}, nil, nil, true},
		{"test1", args{parts, ""}, []string{"UART", "USB"}, []string{"USB", "UART"}, false},
		{"fail4", args{parts1, ""}, nil, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//			t.Parallel()
			got, err := getPeripherals(tt.args.contextMap, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPeripherals() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want1) && !reflect.DeepEqual(got, tt.want2) {
				t.Errorf("getPeripherals() %s = %v, want %v", tt.name, got, tt.want1)
			}
		})
	}
}

func Test_getVirtualMode(t *testing.T) {
	t.Parallel()

	var parts = make(map[string]map[string]string)
	parts["USB"] = map[string]string{"VirtualModexx": "vvmm"}

	var parts1 = make(map[string]map[string]string)
	parts1["USB"] = map[string]string{"xx": "vvmm"}

	type args struct {
		contextMap map[string]map[string]string
		peripheral string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{parts, "USB"}, "vvmm"},
		{"test1", args{parts1, "USB"}, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getVirtualMode(tt.args.contextMap, tt.args.peripheral); got != tt.want {
				t.Errorf("getVirtualMode() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getPins(t *testing.T) {
	var pins = make(map[string]map[string]string)
	pins["PB8"] = map[string]string{"Mode": "MII", "Signal": "I2C1_sss", "GPIO_Label": "lab"}

	pindef := PinDefinition{p: "PB8", pin: "GPIO_PIN_8", port: "GPIOB", mode: "GPIO_MODE_AF_OD", pull: "GPIO_NOPULL", speed: "GPIO_SPEED_FREQ_LOW", alternate: "GPIO_AF4_I2C1"}
	var pindefs = make(map[string]PinDefinition)
	pindefs["I2C1_sss"] = pindef

	type args struct {
		contextMap map[string]map[string]string
		filename   string
		peripheral string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]PinDefinition
		wantErr bool
	}{
		{"test", args{pins, "../../testdata/stm32cubemx/test_msp.c", "I2C1"}, pindefs, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.filename)
			if err != nil {
				t.Errorf("getPinConfiguration() %s cannot open %s", tt.name, tt.args.filename)
				return
			}
			got, err := getPins(tt.args.contextMap, file, tt.args.peripheral)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPins() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPins() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_replaceSpecialChars(t *testing.T) {
	t.Parallel()

	type args struct {
		label string
		ch    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"abcde", "_"}, "abcde"},
		{"test2", args{"a<cde", "_"}, "a_cde"},
		{"test3", args{"ab.,e", "_"}, "ab__e"},
		{"test4", args{"?bcd,", "_"}, "_bcd_"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := replaceSpecialChars(tt.args.label, tt.args.ch); got != tt.want {
				t.Errorf("replaceSpecialChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDigitAtEnd(t *testing.T) {
	t.Parallel()

	type args struct {
		pin string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"abcd12"}, "12"},
		{"test2", args{"abc12e"}, ""},
		{"test3", args{"abcd1"}, "1"},
		{"test4", args{"12"}, "12"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getDigitAtEnd(tt.args.pin); got != tt.want {
				t.Errorf("getDigitAtEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPinConfiguration(t *testing.T) {
	pins := PinDefinition{p: "PB8", pin: "GPIO_PIN_8", port: "GPIOB", mode: "GPIO_MODE_AF_OD", pull: "GPIO_NOPULL", speed: "GPIO_SPEED_FREQ_LOW", alternate: "GPIO_AF4_I2C1"}

	type args struct {
		filename   string
		peripheral string
		pin        string
		label      string
	}
	tests := []struct {
		name    string
		args    args
		want    PinDefinition
		wantErr bool
	}{
		{"test.msp", args{"../../testdata/stm32cubemx/test_msp.c", "", "PB8", ""}, pins, false},
		{"test1.msp", args{"../../testdata/stm32cubemx/test_msp1.c", "", "PB8", ""}, pins, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.filename)
			if err != nil {
				t.Errorf("getPinConfiguration() %s cannot open %s", tt.name, tt.args.filename)
				return
			}
			got, err := getPinConfiguration(file, tt.args.peripheral, tt.args.pin, tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPinConfiguration() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPinConfiguration() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_mxDeviceWriteHeader(t *testing.T) {
	var b bytes.Buffer

	now := time.Now()
	dtString := now.Format("02/01/2006 15:04:05")
	line1 := "/******************************************************************************\n" +
		" * File Name   : fileName\n" +
		" * Date        : " + dtString + "\n" +
		" * Description : STM32Cube MX parameter definitions\n" +
		" * Note        : This file is generated with a generator out of the\n" +
		" *               STM32CubeMX project and its generated files (DO NOT EDIT!)\n" +
		" ******************************************************************************/\n\n" +
		"#ifndef __MX_DEVICE_H\n" +
		"#define __MX_DEVICE_H\n\n"

	type args struct {
		out   *bufio.Writer
		fName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"header", args{fName: "fileName"}, line1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)

			if err := mxDeviceWriteHeader(tt.args.out, tt.args.fName); (err != nil) != tt.wantErr {
				t.Errorf("mxDeviceWriteHeader() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.print() err = %v", err)
			}
			if !errors.Is(err, io.EOF) && str != tt.want {
				t.Errorf("Output.print() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}

func Test_mxDeviceWritePeripheralCfg(t *testing.T) {
	var b bytes.Buffer

	var pins map[string]PinDefinition

	pins1 := make(map[string]PinDefinition)
	pins1["pin1"] = PinDefinition{p: "p"}

	pins2 := make(map[string]PinDefinition)
	pins2["pin2"] = PinDefinition{pin: "pin"}

	pins3 := make(map[string]PinDefinition)
	pins3["pin3"] = PinDefinition{port: "port"}

	pins4 := make(map[string]PinDefinition)
	pins4["pin4"] = PinDefinition{mode: "mode"}

	pins5 := make(map[string]PinDefinition)
	pins5["pin5"] = PinDefinition{pull: "pull"}

	pins6 := make(map[string]PinDefinition)
	pins6["pin6"] = PinDefinition{speed: "speed"}

	pins7 := make(map[string]PinDefinition)
	pins7["pin7"] = PinDefinition{alternate: "alternate"}

	linev := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Virtual mode */\n" +
		"#define MX_peripheral_VM                       vmode\n" +
		"#define MX_peripheral_vmode                    1\n"

	line1 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin1 */\n" +
		"#define MX_pin1_Pin                            p\n"

	line2 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin2 */\n" +
		"#define MX_pin2_GPIO_Pin                       pin\n\n"

	line3 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin2 */\n" +
		"#define MX_pin2_GPIOx                          port\n\n"

	line4 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin3 */\n" +
		"#define MX_pin3_GPIO_Mode                      mode\n\n"

	line5 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin4 */\n" +
		"#define MX_pin4_GPIO_PuPd                      pull\n\n"

	line6 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin5 */\n" +
		"#define MX_pin5_GPIO_Speed                     speed\n"

	line7 := "\n/*------------------------------ peripheral     -----------------------------*/\n" +
		"#define MX_peripheral                          1\n\n" +
		"/* Pins */\n\n" +
		"/* pin6 */\n" +
		"#define MX_pin6_GPIO_AF                        alternate\n\n"

	type args struct {
		out        *bufio.Writer
		peripheral string
		vmode      string
		i2cInfo    map[string]string
		pins       map[string]PinDefinition
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"device", args{peripheral: "peripheral", vmode: "vmode", pins: pins}, linev, false},
		{"device_pin1", args{peripheral: "peripheral", pins: pins1}, line1, false},
		{"device_pin2", args{peripheral: "peripheral", pins: pins2}, line2, false},
		{"device_pin3", args{peripheral: "peripheral", pins: pins3}, line3, false},
		{"device_pin4", args{peripheral: "peripheral", pins: pins4}, line4, false},
		{"device_pin5", args{peripheral: "peripheral", pins: pins5}, line5, false},
		{"device_pin6", args{peripheral: "peripheral", pins: pins6}, line6, false},
		{"device_pin7", args{peripheral: "peripheral", pins: pins7}, line7, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)

			if err := mxDeviceWritePeripheralCfg(tt.args.out, tt.args.peripheral, tt.args.vmode, "", tt.args.i2cInfo, "", "", tt.args.pins); (err != nil) != tt.wantErr {
				t.Errorf("mxDeviceWritePeripheralCfg() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.print() err = %v", err)
			}
			if !errors.Is(err, io.EOF) && str != tt.want {
				t.Errorf("Output.print() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}

func Test_writeDefine(t *testing.T) {
	var b bytes.Buffer

	line1 := "#define name                                    value\n"
	line2 := "#define na_e                                    val_e\n"

	type args struct {
		out   *bufio.Writer
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test", args{name: "name", value: "value"}, line1, false},
		{"test=-", args{name: "na=e", value: "val-e"}, line2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.out = bufio.NewWriter(&b)

			if err := writeDefine(tt.args.out, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("writeDefine() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.args.out.Flush()
			str, err := b.ReadString('\000')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("Output.print() err = %v", err)
			}
			if !errors.Is(err, io.EOF) && str != tt.want {
				t.Errorf("Output.print() %s = %v, want %v", tt.name, str, tt.want)
			}
		})
	}
}
