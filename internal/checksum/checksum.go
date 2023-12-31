/*
 * Copyright 2023 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package checksum

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

// MD5Sum returns the md5 checksum of the given file encoded as a hex string.
func MD5Sum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// separator is used to separate the checksum from the filepath in the checksum file.
const separator = " *"

// CreateChecksumFile creates a checksum file for the given file.
// It returns the path to the created checksum file.
func CreateChecksumFile(file string, checksum string, extension string) (string, error) {
	// Don't allow empty checksum
	if checksum == "" {
		panic("checksum is empty")
	}

	checksumFile := file + extension
	content := checksum + separator + filepath.Base(file) + "\n"

	// Create checksum file
	f, err := os.Create(checksumFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Write content to checksum file
	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	return checksumFile, nil
}
