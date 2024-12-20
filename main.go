package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-ppdiff/domtree"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func readFile(src string) (string, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeFile(t, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(t)
	if err != nil {
		return err
	}
	return nil
}

func execDiff(origin, revised string) error {
	org, err := readFile(origin)
	if err != nil {
		return err
	}
	rev, err := readFile(revised)
	if err != nil {
		return err
	}
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = 0

	diffs := dmp.DiffMain(org, rev, false)
	markup := dmp.DiffPrettyHtml(diffs)

	title := fmt.Sprintf("'%s'→'%s'", filepath.Base(origin), filepath.Base(revised))

	var dt domtree.DomTree
	dt.Init(markup)

	doc := domtree.NewHtmlNode("ja")
	h := domtree.NewHeadNode(title)
	doc.AppendChild(h)
	doc.AppendChild(dt.ToBody(title))
	o := strings.TrimSuffix(origin, filepath.Ext(origin)) + "_diff.html"

	return writeFile(domtree.Decode(doc), o)
}

func run(origin, revised string) int {
	err := execDiff(origin, revised)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func main() {
	var (
		origin  string
		revised string
	)
	flag.StringVar(&origin, "origin", "", "original file path")
	flag.StringVar(&revised, "revised", "", "revised file path")
	flag.Parse()
	os.Exit(run(origin, revised))
}
