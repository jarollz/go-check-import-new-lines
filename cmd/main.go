package main

import (
	"flag"
	"github.com/jarollz/go-check-import-new-lines/internal"
	"log"
)

func main() {
	maxNewLines := flag.Int("n", 1, "maximum number of new lines allowed in golang imports")
	filePath := flag.String("f", "", "the file name to check")
	flag.Parse()

	parser, err := internal.New(int32(*maxNewLines), *filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = parser.ValidateImportsNewLines()
	if err != nil {
		log.Fatal(err)
	}
}
