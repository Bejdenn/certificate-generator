package data

import (
	"strconv"
	"strings"
)

type Participant struct {
	ID         int
	Name, Time string
}

func NewParticipant(ID int, name, time string) *Participant {
	return &Participant{ID: ID, Name: name, Time: time}
}

func (p Participant) FullID(digits int) string {
	numId := strconv.Itoa(p.ID)

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

	return numId + "_" + strings.ReplaceAll(strings.ToLower(p.Name), " ", "")
}
