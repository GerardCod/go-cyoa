package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GerardCod/go-adventure/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web application on")
	file := flag.String("file", "gopher.json", "The json file with the CYOA story")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *file)

	content, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(content)
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New("").Parse(storyTempl))

	h := cyoa.NewHandler(
		story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFunc(pathFunc),
	)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server at: %d.\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	return path[len("/story/"):]
}

var storyTempl = `
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
				<li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
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
