package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/acharyab/gophercises/ex4-html-parser/link"
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

	fmt.Printf("URL links: %+v", links)

}
