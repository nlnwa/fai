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

package log

import (
	"context"
	"log/slog"
)

// noopHandler is a slog.Handler that does nothing.
type noopHandler struct{}

func (n noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (n noopHandler) Handle(context.Context, slog.Record) error { return nil }
func (n noopHandler) WithAttrs([]slog.Attr) slog.Handler        { return n }
func (n noopHandler) WithGroup(string) slog.Handler             { return n }

func Noop() *slog.Logger {
	return slog.New(noopHandler{})
}
