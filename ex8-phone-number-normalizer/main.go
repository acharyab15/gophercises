package main

import "bytes"

func main() {

	number := "123-456-7890"
	normalize(number)
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
