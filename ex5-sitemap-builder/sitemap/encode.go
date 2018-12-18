package sitemap

import (
	"encoding/xml"

	"github.com/acharyab/gophercises/ex4-html-parser/link"
)

// URL represents a URL location
type url struct {
	Loc string `xml:"loc"`
}

// URLSet represents a sitemap
type URLSet struct {
	XMLName    xml.Name `xml:"urlset"`
	SitemapURL []url    `xml:"url"`
}

// GetURLFromLinks extracts a URL
func GetURLFromLinks(links []link.Link) (URLSet, error) {
	uSet := &URLSet{}
	for _, link := range links {
		urlloc := link.Href
		urll := url{Loc: urlloc}
		uSet.SitemapURL = append(uSet.SitemapURL, urll)
	}
	return *uSet, nil
}
