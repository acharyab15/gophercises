package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/acharyab/gophercises/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s. \n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		// not a great idea to panic in general
		panic(err)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)

}
