package scrape

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	o := ParseFromJSON(html, `{
		"title": "head > title",
		"headline": "body h1",
		"sections": {
			"#each": ".container",
			"head": "h2",
			"text": "p",
			"inner": {
				"#each": ".inner",
				"headline": "h3"
			},
			"one": {
				"#element": ".one",
				"value": ".val"
			}
		}
	}`)
	require.Equal(t, o, nil)

	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(b))
}

func TestParser2(t *testing.T) {
	o := ParseFromJSON(html, `{
		"#each": ".container",
		"head": "h2",
		"text": "p"
	}`)

	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(b))
}

var html = `
<html>
	<head>
		<title>Title</title>
	</head>
	<body>
		<h1>Headline</h1>
		<div class="container">
			<h2>Section 1</h2>
			<p>
				Paragraph 1
			</p>
			<div class="one">
				<div class="val">One</div>
			</div>
			<div class="inner">
				<h3>Inner H3</h3>
			</div>
			<div class="inner">
				<h3>Inner H3 next</h3>
			</div>
		</div>
		<div class="container">
			<h2>Section 2</h2>
			<p>
				Paragraph 2
			</p>
			<div class="one"><div class="val">Two</div></div>
			<div class="inner">
				<h3>Inner H3 2</h3>
			</div>
			<div class="inner">
				<h3>Inner H3 2 next</h3>
			</div>
		</div>
	</body>
</html>
`
