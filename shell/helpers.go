package shell

import (
	"strconv"
)

func FormatMsTime(v float64) string {
	suffix := "ms"
	if v > 1000.0 {
		suffix = "S"
		v = v / 1000.0
	}
	return strconv.FormatFloat(v, 'f', calcFloatPrec(v), 64) + suffix
}

func calcFloatPrec(v float64) int {
	prec := 3
	if v >= 1000.000 {
		prec = 0
	} else if v > 0.0 && v < .001 {
		prec = 6
	} else if v > 0.0 && v < 1.0 {
		prec = 4
	} else if v == 0.0 {
		prec = 0
	}
	return prec
}
