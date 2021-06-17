package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseCSV(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(f.Name()) // clean up

	id, name, time := 0, "Max Mustermann", "0:00,0"
	_, _ = f.Write([]byte("Participant, Time"))
	_, _ = f.Write([]byte("\n"))
	_, _ = f.Write([]byte(fmt.Sprintf(`"%q", "%q"`, name, time)))

	res, c := ParseCSV(f)
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	// Count must not be 0 if the passed CSV has one participant
	if c == 0 {
		t.Fatalf("ParseCSV(f) c = %v, want %v", c, 1)
	}

	p := res[0]
	if p.ID != id || p.Name != name || p.Time != time {
		t.Fatalf("ParseCSV(f) = %v %v %v; want %v %v %v", p.ID, p.Name, p.Time, id, name, time)
	}
}

func TestParseCSVRow(t *testing.T) {
	pName := "Max Mustermann"
	pTime := "7:50"
	in := []string{pName, pTime}
	res := parseCSVRow(4, in)
	if res.ID != 4 && res.Name == pName && res.Time == pTime {
		t.Errorf("parseCSVRow(4, in) = Participant with %v %v %v; want %v %v %v", res.ID, res.Name, res.Time, 4, pName, pTime)
	}
}
