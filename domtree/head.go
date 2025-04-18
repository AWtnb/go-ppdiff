package domtree

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

const css = `
#diff-container {
  width: 600px;
  margin: auto;
  font-family: "HackGen";
  font-size: 16px;
  line-height: 1.25;
  word-break: break-all;
}

ins {
  border-radius: 4px;
  background: #ffbebe;
  border: 2px solid tomato;
  text-decoration: none;
}

del+ins {
  border-style: dashed;
  border-width: 1px;
  border-radius: 0 4px 4px 0;
}

del {
  background: #a4e5ff;
  border: 1px solid #05374b;
  color: #929292;
  user-select: none;
`

func getHeadMarkup(title string) string {
	var buf strings.Builder

	buf.WriteString(`<meta charset="utf-8">`)
	buf.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">`)

	faviconHex := "1f4dd"
	fv := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50%%" y="50%%" style="dominant-baseline:central;text-anchor:middle;font-size:90px;">&#x%s;</text></svg>`, faviconHex)
	buf.WriteString(fmt.Sprintf(`<link rel="icon" href="data:image/svg+xml,%s">`, url.PathEscape(fv)))

	buf.WriteString(fmt.Sprintf(`<title>%s</title>`, title))

	return fmt.Sprintf(`<head>%s</head>`, buf.String())
}

func NewHeadNode(title string) *html.Node {
	head := newHeadNode()
	m := getHeadMarkup(title)
	h, _ := html.ParseFragment(strings.NewReader(m), newHeadNode())
	for _, n := range h {
		head.AppendChild(n)
	}

	s := newStyleNode()
	s.AppendChild(newTextNode(css))
	head.AppendChild(s)

	return head
}
