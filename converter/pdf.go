package converter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

const (
	PDF = "PDF"
)

// ConvertError is an error that is raised when some conversion step was not successful.
type ConvertError struct {
	err      error
	destType string
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf("could not convert to %v: %q", e.destType, e.err)
}

// HTMLToPDF converts the passed HTML string to a PDF which contents are returned as a byte slice.
// The HTML string must not be empty, otherwise an error will be returned.
func HTMLToPDF(htmlStr string) ([]byte, error) {
	if len(htmlStr) == 0 {
		return nil, &ConvertError{errors.New("cannot convert to %v with empty html string"), PDF}
	}

	// PDF creation
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, &ConvertError{err, PDF}
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlStr)))

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		return nil, &ConvertError{err, PDF}
	}

	return pdfg.Buffer().Bytes(), nil
}
