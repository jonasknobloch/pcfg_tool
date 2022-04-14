package pcfg

import "strconv"

func FormatWeight(weight float64) string {
	return strconv.FormatFloat(weight, 'f', -1, 64)
}
