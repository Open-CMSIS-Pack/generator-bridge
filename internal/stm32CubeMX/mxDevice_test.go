/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

import (
	"reflect"
	"testing"
)

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

func Test_getContexts(t *testing.T) {
//	t.Parallel()

	var mcuContext1 = make(map[string]map[string]string)
	mcuContext1["Mcu"] = map[string]string{"ContextTest": "myContext"}

	var mcuContext2 = make(map[string]map[string]string)
	mcuContext2["Mcu"] = map[string]string{"Context1": "myContext1", "Context2": "myContext2"}

	var oneEmpty = make(map[int]string)

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
//			t.Parallel()
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

func Test_getDeviceFamily(t *testing.T) {
	t.Parallel()

	var parts = make(map[string]map[string]string)

	parts["Mcu"] = map[string]string{"Family": "STM32U5"}

	var parts1 = make(map[string]map[string]string)

	parts1["Mcu"] = map[string]string{"xFamily": "STM32U5"}

	type args struct {
		contextMap map[string]map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test", args{parts}, "STM32U5", false},
		{"fail", args{parts1}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getDeviceFamily(tt.args.contextMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDeviceFamily() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDeviceFamily() %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getPeripherals1(t *testing.T) {
	t.Parallel()

	var parts = make(map[string]map[string]string)
	parts["Mcu"] = map[string]string{"IP1": "UART", "IP2": "xx", "IP3": "USB3"}

	var parts1 = make(map[string]map[string]string)
	parts1["xxx"] = map[string]string{}

	type args struct {
		contextMap map[string]map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want1   []string
		want2   []string
		wantErr bool
	}{
		{"test", args{parts}, []string{"UART", "USB3"}, []string{"USB3", "UART"}, false},
		{"fail", args{parts1}, nil, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getPeripherals1(tt.args.contextMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPeripherals1() %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want1) && !reflect.DeepEqual(got, tt.want2) {
				t.Errorf("getPeripherals1() %s = %v, want %v", tt.name, got, tt.want1)
			}
		})
	}
}

func Test_getPeripherals(t *testing.T) {
	t.Parallel()

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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func Test_replaceSpecialChars(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceSpecialChars(tt.args.label, tt.args.ch); got != tt.want {
				t.Errorf("replaceSpecialChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDigitAtEnd(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := getDigitAtEnd(tt.args.pin); got != tt.want {
				t.Errorf("getDigitAtEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
