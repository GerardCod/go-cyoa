package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/GerardCod/go-adventure/cyoa"
)

func main() {
	file := flag.String("file", "gopher.json", "The json file with the CYOA story")

	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *file)

	content, err := os.Open(*file)

	if err != nil {
		panic(err)
	}

	var story cyoa.Story

	d := json.NewDecoder(content)

	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
