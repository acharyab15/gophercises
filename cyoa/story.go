package cyoa

import (
	"encoding/json"
	"io"
)

// JSONStory takes in a io.Reader as an input and decodes it into a Story
func JSONStory(r io.Reader) (Story, error) {
	// like marshal/unmarshal for a reader object
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story is a map of string to story Chapter
type Story map[string]Chapter

// Chapter contains a Title, Paragraphs and Options to choose from
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// An Option has a Text and a Chapter
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
