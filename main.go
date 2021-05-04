package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const (
	PORT          = "8080"
	BASE_URL      = "http://localhost:" + PORT
	TEMPLATE_NAME = "web/template/certificate_template"
)

type Participant struct {
	Id   int
	Name string
	Time string
}

func main() {
	argsWithoutProg := os.Args[1:]

	fileName := argsWithoutProg[0]

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("No such file or path: %v", fileName)
	}

	res := parseAll(f)

	// Exposing a web server with individual URLs for every certificate to scrape them later
	mux := http.NewServeMux()

	ids := make([]string, 0)
	for _, v := range res {
		// generate ID from participant name
		id := strings.ReplaceAll(strings.ToLower(v.Name), " ", "") + strconv.Itoa(v.Id)

		mux.Handle("/"+id+"/", NewCert(v))
		ids = append(ids, id)
	}

	// this handler is needed to check if the server is running
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "This is the server for the certificate creator!")
	})

	server := http.Server{Handler: mux, Addr: ":" + PORT}

	go func() {
		defer server.Close()

		for {
			// let the loop run until the server responds
			if !pingServer() {
				continue
			}

			// then start creating the PDF files if everything is up and running
			for _, id := range ids {
				generatePDF(BASE_URL, id, argsWithoutProg[1])
				log.Printf("Created certificate for participant %v\n", id)
			}

			break
		}
	}()

	log.Fatal(server.ListenAndServe())
}

func pingServer() bool {
	time.Sleep(time.Second)

	log.Println("Waiting for pages to be served...")
	resp, err := http.Get(BASE_URL)
	if err != nil {
		log.Println("Failed: ", err)
		return false
	}

	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Not OK: ", resp.StatusCode)
		return false
	}

	return true
}

func generatePDF(baseUrl string, id string, outputPath string) {
	os.Mkdir(outputPath, os.ModePerm)
	// PDF creation
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalln("Could not create PDF generator", err)
	}

	url := baseUrl + "/" + id + "/"
	htmlStr := fetch(url)

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

func fetch(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

type Certificate struct {
	v Participant
}

func NewCert(v Participant) Certificate {
	return Certificate{v: v}
}

var templates = template.Must(template.ParseFiles(TEMPLATE_NAME + ".html"))

func (c Certificate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(rw, TEMPLATE_NAME+".html", c.v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func parseAll(f *os.File) []Participant {
	r := csv.NewReader(f)

	// just read the whole file
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Error while reading CSV")
	}

	res := make([]Participant, 0)

	// skip the first entry because this will be the CSV header
	for i, v := range records[1:] {
		res = append(res, parse(i, v))
	}

	return res
}

func parse(line int, in []string) Participant {
	if len(in) == 0 {
		log.Fatalln("Tried to parse CSV entry that had no values")
	}

	return Participant{Id: line, Name: in[0], Time: in[1]}
}
