package main

/*
	Refresh the static.go file with the latest equivalent  pairs from an
 	EquivalentUnifiedIdeograph.txt file. Run this command from the root
	of the repo:` `go run cmd/updatestatic/main.go``
*/

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitranim/repr"

	eqi "github.com/mochi-co/equivalent-unified-ideograph"
)

func init() {
	index = flag.String("index", "data/EquivalentUnifiedIdeograph.txt", "the latest version of EquivalentUnifiedIdeograph.txt")
	output = flag.String("output", "static.go", "the static go file containing the equivalent pairs")
	flag.Parse()
}

var (
	index  *string
	output *string
)

func main() {
	file, err := os.Open(*index)
	if err != nil {
		log.Fatalln("Couldn't open the index file", err)
	}

	pairs, err := eqi.ExtractPairs(file)
	if err != nil {
		log.Fatalln(err)
	}

	d := `package eqi 

var Pairs = ` + repr.String(pairs)

	d = strings.Replace(d, "eqi.", "", -1)
	err = ioutil.WriteFile(*output, []byte(d), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Sucessfully updated %s with data from %s", *output, *index)
}
