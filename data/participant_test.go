package data

import (
	"regexp"
	"testing"
)

func TestFullID(t *testing.T) {
	pId, pName, time := 123, "maxmustermann", "" // the time value is irrelevant for this test case
	p := NewParticipant(pId, pName, time)
	res := p.FullID(4)

	want := regexp.MustCompile(`0123_maxmustermann`)
	if !want.MatchString(res) {
		t.Fatalf(`p.FullID(4) = %q, want match for %#q`, res, want)
	}
}
