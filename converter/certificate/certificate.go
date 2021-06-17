package certificate

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/Bejdenn/certificate-generator/converter"
	"github.com/Bejdenn/certificate-generator/parser"
)

func GenerateAll(f *os.File, destPath string, t *template.Template) {
	res, count := parser.ParseCSV(f)
	err := os.Mkdir(destPath, os.ModePerm)
	if err != nil {
		return
	}

	for _, v := range res {
		// generate ID from participant name
		maxPad := len(strconv.Itoa(count))
		id := v.FullID(maxPad)

		var buf bytes.Buffer

		err := t.Execute(&buf, v)
		if err != nil {
			log.Printf("Could not execute template: %v", err)
			continue
		}

		pdf, err := converter.HTMLToPDF(buf.String())
		if err != nil {
			log.Fatalf("Could not create certificates: %q", err)
		}

		err = ioutil.WriteFile(destPath+"/"+id+".pdf", pdf, 0666)
		if err != nil {
			log.Fatalf("Could not create certificates: %q", err)
		}

		log.Printf("Created certificate for participant %v\n", id)
	}
}
