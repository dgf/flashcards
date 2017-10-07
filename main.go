package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/dgf/flashcards/dictionary"
	"github.com/dgf/flashcards/pdf"
)

var (
	config = flag.String("config", "config.toml", "configuration file")
	dict   = flag.String("dict", "dicts/es-de.dict", "card dictionary file")
	path   = flag.String("path", "out.pdf", "output PDF file")
)

func init() {
	flag.Parse()
}

func main() {

	log.Printf("decode configuration: %s", *config)
	var options pdf.Options
	if _, err := toml.DecodeFile(*config, &options); err != nil {
		log.Fatal(err)
	}
	log.Printf("options: %v", options)

	log.Printf("read dictionary: %s", *dict)
	cards, err := dictionary.Read(*dict)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d cards", len(cards))

	log.Printf("print PDF: %s", *path)
	file := pdf.New(*path, options)
	if err := file.Print(cards); err != nil {
		log.Fatal(err)
	}
}
