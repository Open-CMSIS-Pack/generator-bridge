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

func TestAddLine(t *testing.T) {
	var text utils.TextBuilder
	text.AddLine("A line")
	assert.Equal(t, "A line\n", text.GetLine())
	text.AddLine("A second line")
	assert.Equal(t, "A line\nA second line\n", text.GetLine())
}
