/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"testing"

	"github.com/open-cmsis-pack/generator-bridge/internal/utils"
	"github.com/stretchr/testify/assert"
)

func testAddLine(t *testing.T, pid int) {
	var text utils.TextBuilder
	text.AddLine("A line")
	assert.Equal(t, text.GetLine(), "A line")
	text.AddLine("A second line")
	assert.Equal(t, text.GetLine(), "A line\nA second line")
}
