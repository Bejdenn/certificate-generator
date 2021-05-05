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
	port             = "8080"
	baseURL          = "http://localhost:" + port
	certTemplateName = "./web/template/certificate_template"
)

type Participant struct {
	ID         int
	Name, Time string
}

func main() {
	argsWithoutProg := os.Args[1:]

	fileName := argsWithoutProg[0]

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("No such file or path: %v", fileName)
	}

	res, count := parseAll(f)

	// Exposing a web server with individual URLs for every certificate to scrape them later
	mux := http.NewServeMux()

	ids := make([]string, 0)
	for _, v := range res {
		// generate ID from participant name
		id := createCertID(v.ID, v.Name, len(strconv.Itoa(count)))
		mux.Handle("/"+id+"/", NewCert(v))
		ids = append(ids, id)
	}

	// this handler is needed to check if the server is running
	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "This is the server for the certificate creator!")
	})

	server := http.Server{Handler: mux, Addr: ":" + port}

	go func() {
		defer server.Close()

		for {
			// let the loop run until the server responds
			if !pingServer() {
				continue
			}

			// then start creating the PDF files if everything is up and running
			for _, id := range ids {
				generatePDF(baseURL, id, argsWithoutProg[1])
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
	resp, err := http.Get(baseURL)
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

var templates = template.Must(template.New("certificate_template").Parse(TemplateStr))

func (c Certificate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := templates.Execute(rw, c.v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

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
