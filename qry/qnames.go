package qry

import (
	"math/rand"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// Generating Preductable query name
func PQname(labels int, suffix int) string {
	var label string
	for i := 1; i <= (labels - 1); i++ {
		labelprefix := "foo"
		labelfmt := labelprefix + strconv.Itoa(suffix) + "."
		label = labelfmt + label
	}
	label = label + "lab"
	return label
}

//Random string
func Rstring(l int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//Generating Random query name
func RQname(labels int) string {
	var label string
	for i := 1; i <= (labels - 1); i++ {
		labelfmt := Rstring(16) + "."
		label = labelfmt + label
	}
	label = label + "lab"
	return label
}
