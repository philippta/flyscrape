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
