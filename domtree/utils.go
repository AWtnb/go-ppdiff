package domtree

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func elemAttr(name, val string) html.Attribute {
	return html.Attribute{Key: name, Val: val}
}

func appendAttr(node *html.Node, name, val string) {
	node.Attr = append(node.Attr, elemAttr(name, val))
}

func removeAttr(node *html.Node, name string) {
	attrs := []html.Attribute{}
	for _, attr := range node.Attr {
		if attr.Key != name {
			attrs = append(attrs, attr)
		}
	}
	node.Attr = attrs
}

func setId(node *html.Node, id string) {
	removeAttr(node, "id")
	appendAttr(node, "id", id)
}

func newTextNode(data string) *html.Node {
	n := &html.Node{
		Type: html.TextNode,
		Data: data,
	}
	return n
}

func newElementNode(data string, dataAtom atom.Atom) *html.Node {
	n := &html.Node{
		Type:     html.ElementNode,
		Data:     data,
		DataAtom: dataAtom,
	}
	return n
}

func newHeadNode() *html.Node {
	return newElementNode("head", atom.Head)
}

func newDivNode() *html.Node {
	return newElementNode("div", atom.Div)
}

func newH1Node(t string) *html.Node {
	n := newElementNode("h1", atom.H1)
	n.AppendChild(newTextNode(t))
	return n
}

func newStyleNode() *html.Node {
	n := newElementNode("style", atom.Style)
	return n
}

func NewHtmlNode(lang string) *html.Node {
	n := newElementNode("html", atom.Html)
	appendAttr(n, "lang", lang)
	return n
}

func Decode(node *html.Node) string {
	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>")
	html.Render(&buf, node)
	return strings.ReplaceAll(buf.String(), "\u00B6", "\u21B5")
}
