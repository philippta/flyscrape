package scrape

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	html := `<html>
<body>
	<h1 id="title">Page Title</h1>
	<div id="posts">
		<div class="post">First post</div>
		<div class="post">Second post</div>
		<div class="post">Third post</div>
	</div>
</body>
</html>`

	title := Query(Doc(html), "#title")
	require.Equal(t, title, "Page Title")
}
