package utility

import (
	"os"
	"strconv"
)

func FormatWeight(weight float64) string {
	return strconv.FormatFloat(weight, 'f', -1, 64)
}

func OpenFile(name string) (*os.File, error) {
	if name == "" {
		return os.Stdin, nil
	}

	return os.Open(name)
}

func CreateFile(name string) (*os.File, error) {
	if name == "" {
		return os.Stdout, nil
	}

	return os.Create(name)
}
