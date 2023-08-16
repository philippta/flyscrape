package flyscrape

import (
	"encoding/json"
	"os"
)

func PrettyPrint(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "   ")
	enc.Encode(v)
}
