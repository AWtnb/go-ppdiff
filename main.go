/*
Copyright 2024 AWtnb

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

func toLineFeed(s string) string {
	c := "\n"
	return strings.NewReplacer(
		"\r\n", c,
		"\r", c,
	).Replace(s)
}

func readFile(src string) (string, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	return toLineFeed(string(b)), nil
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

func toStem(p string) string {
	return strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
}

func execDiff(origin, revised, out string) error {
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
	dmp.DiffCleanupSemanticLossless(diffs)

	title := fmt.Sprintf("'%s'→'%s'", filepath.Base(origin), filepath.Base(revised))

	var dt domtree.DomTree
	dt.Init(diffs)

	doc := domtree.NewHtmlNode("ja")
	h := domtree.NewHeadNode(title)
	doc.AppendChild(h)
	doc.AppendChild(dt.ToBody(title))

	if len(out) < 1 {
		n := fmt.Sprintf("%s_diff_from_%s.html", toStem(revised), toStem(origin))
		out = filepath.Join(filepath.Dir(revised), n)
	}
	if !strings.HasSuffix(out, ".html") {
		out = out + ".html"
	}

	return writeFile(domtree.Decode(doc), out)
}

func run(origin, revised, out string) int {
	err := execDiff(origin, revised, out)
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
		out     string
	)
	flag.StringVar(&origin, "origin", "", "original file path")
	flag.StringVar(&revised, "revised", "", "revised file path")
	flag.StringVar(&out, "out", "", "output file path")
	flag.Parse()
	os.Exit(run(origin, revised, out))
}
