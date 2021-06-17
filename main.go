package main

import (
	"errors"
	"html/template"
	"log"
	"os"

	"github.com/Bejdenn/certificate-generator/converter/certificate"
)

var t = template.Must(template.New("certificate_template").Parse(TemplateStr))

var missingArgsError = errors.New("missing arguments")

func main() {
	log.Println("Starting certificate generator...")
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 2 {
		panic(missingArgsError)
	}

	fileName := argsWithoutProg[0]

	// read CSV
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	certificate.GenerateAll(f, argsWithoutProg[1], t)
	log.Println("Finished generating certificates!")
}
