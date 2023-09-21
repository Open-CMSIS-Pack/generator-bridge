/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	uuidWithHyphen, _ := uuid.NewRandom()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	fmt.Println(uuid)
	fmt.Println("Hello, World!")
}
