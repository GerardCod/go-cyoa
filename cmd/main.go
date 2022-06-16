package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server at: %d.\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
