package domtree

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// https://github.com/sergi/go-diff/blob/5b0b94c5c0d3d261e044521f7f46479ef869cf76/diffmatchpatch/diff.go#L1119
func toHtml(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	i := 1
	br := `<span class="break"></span><br>`
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", br, -1)
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff.WriteString(fmt.Sprintf("<ins tabindex=\"%d\">", i))
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</ins>")
			i += 1
		case diffmatchpatch.DiffDelete:
			_, _ = buff.WriteString(fmt.Sprintf("<del inert tabindex=\"%d\">", i))
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</del>")
			i += 1
		case diffmatchpatch.DiffEqual:
			_, _ = buff.WriteString("<span>")
			_, _ = buff.WriteString(text)
			_, _ = buff.WriteString("</span>")
		}
	}
	return buff.String()
}

type DomTree struct {
	root *html.Node
}

func (dt *DomTree) Init(diffs []diffmatchpatch.Diff) error {
	markup := toHtml(diffs)
	nodes, err := html.ParseFragment(strings.NewReader(markup), newDivNode())
	if err != nil {
		return err
	}

	d := newDivNode()
	setId(d, "diff-container")

	for _, n := range nodes {
		d.AppendChild(n)
	}
	dt.root = d

	return nil
}

func (dt *DomTree) ToBody(heading string) *html.Node {
	b := newElementNode("body", atom.Body)
	h1 := newH1Node(heading)
	b.AppendChild(h1)
	b.AppendChild(dt.root)
	return b
}
