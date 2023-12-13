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
	"errors"
	"io"

	"github.com/nlnwa/gowarc"
)

func IsValid(file, tmpDir string) (bool, error) {
	wf, err := gowarc.NewWarcFileReader(file, 0,
		gowarc.WithBufferTmpDir(tmpDir),
	)
	if err != nil {
		return false, err
	}
	defer wf.Close()

	for {
		wr, _, validation, err := wf.Next()
		if errors.Is(err, io.EOF) {
			return validation.Valid(), nil
		}
		if err != nil {
			if wr != nil {
				defer wr.Close()
			}
			return false, nil
		}
		func() {
			defer wr.Close()
			err = wr.ValidateDigest(validation)
			if err != nil {
				*validation = append(*validation, err)
			}
		}()
		// stop processing if we have found an invalid record
		if !validation.Valid() {
			return false, nil
		}
	}
}
