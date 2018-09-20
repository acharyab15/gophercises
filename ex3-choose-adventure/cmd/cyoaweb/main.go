package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/acharyab/gophercises/ex3-choose-adventure/cyoa"
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

	// like marshal/unmarshal for a reader object
	d := json.NewDecoder(f)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)

}
