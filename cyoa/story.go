package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerImpl))
}

var defaultHandlerImpl = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose your own adventure</title>
</head>
<body>
	<section class="page">
		<h1>{{.Title}}</h1>
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
			text-align: center;
			position: relative;
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
</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	var chapter Chapter
	var ok bool

	if chapter, ok = h.s[path]; !ok {
		http.Error(w, "Chapter not found", http.StatusNotFound)
	}

	err := tpl.Execute(w, chapter)

	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
	}

}

func JsonStory(r io.Reader) (Story, error) {
	var story Story

	d := json.NewDecoder(r)

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

// Story represents the content for the cyoa content
type Story map[string]Chapter

// Chapter represents a Chapter
type Chapter struct {
	Title      string   `json:"title,omitempty"`
	Paragraphs []string `json:"story,omitempty"`
	Options    []Option `json:"options,omitempty"`
}

type Option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

type Demo struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}
