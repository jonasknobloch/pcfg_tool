package utility

import (
	"io"
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

func CopyToTemp(file *os.File) (*os.File, error) {
	temp, err := os.CreateTemp("", "*")

	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(temp, file); err != nil {
		return nil, err
	}

	return temp, nil
}
