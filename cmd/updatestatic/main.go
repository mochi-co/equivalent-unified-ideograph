package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	eqi "github.com/mochi-co/equivalent-unified-ideograph"
)

func init() {
	index = flag.String("index", "data/EquivalentUnifiedIdeograph.txt", "the latest version of EquivalentUnifiedIdeograph.txt")
	input = flag.String("input", "input.txt", "the input file containing the unicode characters")
	output = flag.String("output", "output.txt", "the output file to write with the replaced characters")
	flag.Parse()
}

var (
	index  *string
	input  *string
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

	fmt.Println(pairs)
}
