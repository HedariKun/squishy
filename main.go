package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {

	inputName := flag.String("input", "", "use to define the name of the input file")
	outputNamef := flag.String("output", "", "ues to define the name of the desired output file, by default it will be the same as the input")
	help := flag.Bool("help", false, "show this message")

	flag.Parse()

	args := flag.Args()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if len(args) < 0 && *inputName == "" {
		log.Fatal("you need to specify a sqlite database")
	}

	databaseName := "/"
	outputName := "/"

	if *inputName == "" {
		databaseName += args[0]
	} else {
		databaseName += *inputName
	}

	if *outputNamef == "" {
		outputName += path.Base(databaseName)
	} else {
		outputName += *outputNamef
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(dir)
	}

	println("process started")
	excelizing(dir+databaseName, dir+outputName)
	println("process finished")
}
