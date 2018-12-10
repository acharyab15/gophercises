package link

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Link representes a link (<a href="...">) in an HTML document
type Link struct {
	Href string
	Text string
}

var r io.Reader

// Parse will take in an HTML document and will return a
// slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	dfs(doc, "")

	return nil, nil
}

func dfs(n *html.Node, padding string) {
	msg := n.Data
	if n.Type == html.ElementNode && msg == "a" {
		msg = "<" + msg + ">"
		fmt.Println(padding, msg)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dfs(c, padding+"  ")
	}
}
