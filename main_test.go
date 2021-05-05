package main

import (
	"regexp"
	"testing"
)

func TestCreateCertID(t *testing.T) {
	pId := 123
	pName := "maxmustermann"
	want := regexp.MustCompile(`0123_maxmustermann`)
	res := createCertID(pId, pName, 4)
	if !want.MatchString(res) {
		t.Fatalf(`createCertID(pId, pName, 4) = %q, want match for %#q`, res, want)
	}
}
