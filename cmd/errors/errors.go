/*
 * Copyright (c) 2023 Arm Limited. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package errors

import (
	"errors"
)

// Is returns true if err is equals to target
func Is(err, target error) bool {
	return err == target
}

var (
	// Errors related to file system
	ErrFailedCreatingFile        = errors.New("failed to create a local file")
	ErrFailedWrittingToLocalFile = errors.New("failed writing HTTP stream to local file")
	ErrFailedDecompressingFile   = errors.New("fail to decompress file")
	ErrFailedInflatingFile       = errors.New("fail to inflate file")
	ErrFailedCreatingDirectory   = errors.New("fail to create directory")
	ErrFileNotFound              = errors.New("file not found")
	ErrDirectoryNotFound         = errors.New("directory not found")
	ErrPathAlreadyExists         = errors.New("path already exists")
	ErrCopyingEqualPaths         = errors.New("failed copying files: source is the same as destination")
	ErrMovingEqualPaths          = errors.New("failed moving files: source is the same as destination")

	// Security errors
	ErrInsecureZipFileName = errors.New("zip file contains insecure characters: ../")
	ErrFileTooBig          = errors.New("files cannot be over 20G")
	ErrIndexPathNotSafe    = errors.New("index url path does not start with HTTPS")

	// Cmdline errors
	ErrIncorrectCmdArgs = errors.New("incorrect setup of command line arguments")

	// Error/Flag to detect when a user has requested early termination
	ErrTerminatedByUser = errors.New("terminated by user request")
)
