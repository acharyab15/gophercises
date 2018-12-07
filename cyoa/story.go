package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Choose Your Own Adventure</title>
	</head>
	<body>
		<h1> {{.Title}}</h1>
		{{range .Paragraphs}}
			<p>{{.}}</p>
		{{end}}
		<ul>
			{{range .Options}}
			<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
			{{end}}
		</ul>
	</body>
</html>`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}

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
