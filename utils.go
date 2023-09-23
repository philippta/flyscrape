// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape

import (
	"bytes"
	"encoding/json"
	"strings"
)

func PrettyPrint(v any, prefix string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, "  ")
	enc.Encode(v)
	return prefix + strings.TrimSuffix(buf.String(), "\n")
}

func Print(v any, prefix string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	return prefix + strings.TrimSuffix(buf.String(), "\n")
}

func ParseOptions(opts Options, v any) {
	json.Unmarshal(opts, v)
}
