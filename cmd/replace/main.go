package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	eqi "github.com/mochi-co/equivalent-unified-ideograph"
)

var (
	input  *string
	output *string
)

func init() {
	input = flag.String("input", "example/mismatched.txt", "the input file containing the mismatched ideographs")
	output = flag.String("output", "-", "the output file you wish to write to - omit to overwrite input")
	flag.Parse()
}

func main() {
	file, err := os.Open(*input)
	if err != nil {
		log.Fatalln("couldn't open the input file", err)
	}

	d, err := eqi.BufferedReplace(file)
	if err != nil {
		log.Fatalln("couldn't replace characters", err)
	}

	if *output == "-" {
		output = input
	}

	err = ioutil.WriteFile(*output, d.Bytes(), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
