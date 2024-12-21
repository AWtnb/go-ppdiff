package domtree

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DomTree struct {
	root *html.Node
}

func (dt *DomTree) Init(markup string) error {
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

func (dt *DomTree) cleanStyle() {
	var dfs func(*html.Node)
	dfs = func(node *html.Node) {
		if node.Type == html.ElementNode && node.FirstChild != nil {
			if node.Data == "ins" || node.Data == "del" {
				removeAttr(node, "style")
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(dt.root)
}

func (dt *DomTree) setTabIndex() {
	idx := 1
	var dfs func(*html.Node)
	dfs = func(node *html.Node) {
		if node.Type == html.ElementNode && node.FirstChild != nil {
			if node.Data == "ins" || node.Data == "del" {
				removeAttr(node, "style")
				appendAttr(node, "tabindex", fmt.Sprint(idx))
				idx += 1
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(dt.root)
}

func (dt *DomTree) ToBody(heading string) *html.Node {
	dt.cleanStyle()
	dt.setTabIndex()
	b := newElementNode("body", atom.Body)
	h1 := newH1Node(heading)
	b.AppendChild(h1)
	b.AppendChild(dt.root)
	return b
}
