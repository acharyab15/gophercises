package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/acharyab/gophercises/ex4-html-parser/link"
	"github.com/acharyab/gophercises/ex5-sitemap-builder/sitemap"
)

func main() {
	urlName := flag.String("url", "", "a url to build a sitemap for")
	flag.Parse()

	fmt.Println("The URL is: ", *urlName)
	resp, err := http.Get(*urlName)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	links, err := link.Parse(resp.Body)

	urls, err := sitemap.GetURLFromLinks(links)

	output, err := xml.MarshalIndent(urls, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)

}
