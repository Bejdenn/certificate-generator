package parser

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/Bejdenn/certificate-generator/data"
)

func ParseCSV(f *os.File) (res []*data.Participant, count int) {
	r := csv.NewReader(f)

	// just read the whole file
	records := make([][]string, 0)
	for {
		r, err := r.Read()
		if err == io.EOF {
			break
		}

		records = append(records, r)

	}

	res = make([]*data.Participant, 0)

	// skip the first entry because this will be the CSV header
	for i, v := range records[1:] {
		res = append(res, parseCSVRow(i, v))
	}

	return res, len(records[1:])
}

func parseCSVRow(lineNr int, in []string) *data.Participant {
	if len(in) == 0 {
		log.Fatalln("Tried to parse CSV entry that had no values")
	}

	return data.NewParticipant(lineNr, in[0], in[1])
}
