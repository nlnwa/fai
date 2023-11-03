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

package warc

import (
	"testing"
)

const (
	validWarcFile   = "testdata/valid.warc.gz"
	invalidWarcFile = "testdata/invalid.warc.gz"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		file    string
		isValid bool
	}{
		{
			file:    validWarcFile,
			isValid: true,
		},
		{
			file:    invalidWarcFile,
			isValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.file, func(t *testing.T) {
			isValid, err := IsValid(test.file, "")
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if isValid != test.isValid {
				t.Errorf("expected %v, got: %v", test.isValid, isValid)
			}
		})
	}
}
