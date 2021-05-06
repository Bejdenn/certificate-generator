package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const (
	port             = "8080"
	baseURL          = "http://localhost:" + port
	certTemplateName = "./web/template/certificate_template"
)

type MissingArgsError struct {
	missingArgs []string
}

func (e MissingArgsError) Error() string {
	return fmt.Sprintf("certificate-generator needs the args %q to be able run", e.missingArgs)
}

func main() {
	log.Println("Starting certificate generator...")
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		panic(MissingArgsError{missingArgs: []string{"CSV file path", "Destination folder path"}})
	}

	fileName := argsWithoutProg[0]

	// read CSV
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	generateCerts(f, argsWithoutProg[1])
	log.Println("Finished generating certificates!")
}

type Participant struct {
	ID         int
	Name, Time string
}

func generateCerts(f *os.File, destPath string) {
	res, count := parseAll(f)
	os.Mkdir(destPath, os.ModePerm)

	for _, v := range res {
		// generate ID from participant name
		id := createCertID(v.ID, v.Name, len(strconv.Itoa(count)))

		var buf bytes.Buffer

		err := templates.Execute(&buf, v)
		if err != nil {
			log.Printf("Could not execute template: %v", err)
			continue
		}

		generatePDF(buf.String(), id, destPath)

		log.Printf("Created certificate for participant %v\n", id)
	}
}

func generatePDF(htmlStr string, id string, outputPath string) {
	// PDF creation
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalln("Could not create PDF generator", err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlStr)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(outputPath + "/" + id + ".pdf")
	if err != nil {
		log.Fatal(err)
	}
}

var templates = template.Must(template.New("certificate_template").Parse(TemplateStr))

func parseAll(f *os.File) (res []Participant, count int) {
	r := csv.NewReader(f)

	// just read the whole file
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Error while reading CSV")
	}

	res = make([]Participant, 0)

	// skip the first entry because this will be the CSV header
	for i, v := range records[1:] {
		res = append(res, parse(i, v))
	}

	return res, len(records[1:])
}

func parse(line int, in []string) Participant {
	if len(in) == 0 {
		log.Fatalln("Tried to parse CSV entry that had no values")
	}

	return Participant{ID: line, Name: in[0], Time: in[1]}
}

func createCertID(id int, name string, digits int) string {
	numId := strconv.Itoa(id)

	// skip padding part if the num ID has enough digits
	if pad := digits - len(numId); pad > 0 {
		var zeros string

		// append as many zeros as padding is needed to fulfill digits count
		for i := 0; i < pad; i++ {
			zeros += "0"
		}

		// new num ID consist of the padding and the participant ID
		numId = zeros + numId
	}

	return numId + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "")
}
