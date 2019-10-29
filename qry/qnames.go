package qry

import "strconv"

func PQname(labels int, suffix int) string {
	var label string
	for i := 1; i <= (labels - 1); i++ {
		labelprefix := "foo"
		labelfmt := labelprefix + strconv.Itoa(suffix) + "."
		label = labelfmt + label
	}
	label = label + "com"
	return label
}
