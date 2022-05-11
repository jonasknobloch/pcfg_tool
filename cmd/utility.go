package cmd

import (
	"log"
	"os"
)

func OpenStdin() *os.File {
	name := os.Getenv("STDIN")

	if name == "" {
		return os.Stdin
	}

	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	return f
}
