package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/acharyab/gophercises/ex4-html-parser/link"
)

/*
1. GET the webpage
2. parse all the links on the page (link package)
3. build proper urls with links
4. filter out any links w/ a diff domain
5. find all the pages (BFS)
6. print out XML
*/

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "a url to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXML := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}

	fmt.Println(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()
}

func bfs(urlStr string, maxDepth int) []string {
	// struct uses less memory than a bool
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}
	var ret []string
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(body io.Reader, base string) []string {
	links, _ := link.Parse(body)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}

}
