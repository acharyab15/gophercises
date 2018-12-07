package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
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
	    <section class="page">
			<h1> {{.Title}}</h1>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}
			<ul>
				{{range .Options}}
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</section>
		<style>
          body {
            font-family: helvetica, arial;
          }
          h1 {
            text-align:center;
            position:relative;
          }
          .page {
            width: 80%;
            max-width: 500px;
            margin: auto;
            margin-top: 40px;
            margin-bottom: 40px;
            padding: 80px;
            background: #FFFCF6;
            border: 1px solid #eee;
            box-shadow: 0 10px 6px -6px #777;
          }
          ul {
            border-top: 1px dotted #ccc;
            padding: 10px 0 0 0;
            -webkit-padding-start: 0;
          }
          li {
            padding-top: 10px;
          }
          a,
          a:visited {
            text-decoration: none;
            color: #6295b5;
          }
          a:active,
          a:hover {
            color: #7792a2;
          }
          p {
            text-indent: 1em;
          }
        </style>
	</body>
</html>`

// HandlerOption is a function which is a pointer to a handler
type HandlerOption func(h *handler)

// WithTemplate takes in a template.Template and returns a HandlerOption
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// NewHandler takes in a Story and HandlerOptions and returns a http.Handler
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	// removing the leading slash
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
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
