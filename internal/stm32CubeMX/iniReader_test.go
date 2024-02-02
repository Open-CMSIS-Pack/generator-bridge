/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stm32cubemx

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
