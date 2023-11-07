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
	"os"
	"path/filepath"
	"testing"
)

func TestMD5Sum(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.CreateTemp(tmpDir, "test")

	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	_, err = f.WriteString("Test string")
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}
	testFile := f.Name()
	defer os.Remove(testFile)
	defer f.Close()

	got, err := MD5Sum(testFile)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	want := "0fd3dbec9730101bff92acc820befc34"
	if got != want {
		t.Errorf("expected %s, but got %s", want, got)
	}
}

func TestCreateChecksumFile(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.CreateTemp(tmpDir, "test")

	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	_, err = f.WriteString("Test string")
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}
	testFile := f.Name()
	defer os.Remove(testFile)
	defer f.Close()

	wantHash := "0fd3dbec9730101bff92acc820befc34"
	fileHash, err := MD5Sum(testFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	checksumFile, err := CreateChecksumFile(testFile, fileHash, ".md5")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer os.Remove(checksumFile)

	b, err := os.ReadFile(checksumFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := string(b)
	want := wantHash + separator + filepath.Base(testFile) + "\n"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}
