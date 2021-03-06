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

// Esta función inicializa los valores de algunas variables al momento de cargar
// el módulo actual en tiempo de ejecución.
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

type HandlerOption func(h *handler)

// WithTemplate es una opción funcional para establecer la plantilla
// que mostrará los datos de la historia.
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFunc es una opción funcional para definir el procesamiento de
// las peticiones.
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// NewHandler es una función para crear una variable de http.Handler
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

// defaultPathFn define la ruta por defecto para el funcionamiento del enrutador.
func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

// ServeHTTP procesa las peticiones del servidor.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	var chapter Chapter
	var ok bool

	if chapter, ok = h.s[path]; !ok {
		http.Error(w, "Chapter not found", http.StatusNotFound)
	}

	err := h.t.Execute(w, chapter)

	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
	}

}

// JsonStory realiza la conversión de un json a una variable de tipo Story.
// Recibe un parámetro de tipo io.Reader para procesar los bytes del json.
// Devuelve ya sea una variable de tipo Story o un error.
func JsonStory(r io.Reader) (Story, error) {
	var story Story

	d := json.NewDecoder(r)

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

// Story es una colección de capítulos para la historia de la aplicación.
type Story map[string]Chapter

// Chapter es una struct que representa un capítulo de la historia.
type Chapter struct {
	Title      string   `json:"title,omitempty"`
	Paragraphs []string `json:"story,omitempty"`
	Options    []Option `json:"options,omitempty"`
}

// Option es uno de los caminos alternos para la historia de la aplicación.
type Option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

// Demo
type Demo struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}
