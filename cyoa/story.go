package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
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
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
    <p>{{.}}</p>
  {{end}}

  <ul>
    {{range .Options}}
      <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
  </ul>
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
	err := tpl.Execute(w, h.s["intro"])

	if err != nil {
		panic(err)
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
