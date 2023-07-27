package scrape

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseFromJSON(html, input string) any {
	var inputJSON map[string]any
	json.Unmarshal([]byte(input), &inputJSON)
	return Parse(html, inputJSON)
}

func Parse(html string, fields map[string]any) any {
	return queryMap(Doc(html), fields)
}

func AddMeta(result any, key string, value any) {
	switch res := result.(type) {
	case []map[string]any:
		for i := range res {
			res[i][key] = value
		}
	case map[string]any:
		res[key] = value
	}
}

func walk(s *goquery.Selection, fields map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range fields {
		if strings.HasPrefix(k, "#") {
			continue
		}

		switch val := v.(type) {
		case string:
			segs := strings.SplitN(k, "#", 2)
			if len(segs) == 2 && segs[1] == "html" {
				out[segs[0]] = QueryHTML(s, val)
			} else if len(segs) == 2 {
				out[segs[0]] = QueryAttr(s, val, segs[1])
			} else {
				out[k] = Query(s, val)
			}

		case map[string]any:
			out[k] = queryMap(s, val)
		}
	}
	return out
}

func queryMap(s *goquery.Selection, fields map[string]any) any {
	if sel, ok := fields["#each"].(string); ok {
		rows := []map[string]any{}
		QueryFunc(s, sel, func(s *goquery.Selection) {
			rows = append(rows, walk(s, fields))
		})
		return rows
	}

	if sel, ok := fields["#element"].(string); ok {
		return walk(s.Find(sel), fields)
	}

	return walk(s, fields)
}
