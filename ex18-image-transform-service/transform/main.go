package main

import (
	"io"
	"os"

	"github.com/acharyab/gophercises/ex18-image-transform-service/transform/primitive"
)

func main() {
	inFile, err := os.Open("abc.jpg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 10)
	if err != nil {
		panic(err)
	}
	os.Remove("out.png")
	outFile, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	io.Copy(outFile, out)
	// fmt.Println(string(out))
}
