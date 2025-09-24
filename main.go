package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"cutc/cutc"
)

var env = "development"
var version = "0.0.0"
var help = fmt.Sprintf(Help, version)

func main() {
	var args = cutc.Args{}
	flag.StringVar(&args.Delimiter, "d", ",", "Fields delimiter")
	flag.StringVar(&args.FieldsList, "f", "", "Fields indexes to cut (starting from 1, order matters)")
	flag.BoolVar(&args.SkipHeader, "h", false, "Skip csv header")
	flag.StringVar(&args.Delimiter, "delimiter", ",", "Alias for -d")
	flag.StringVar(&args.FieldsList, "fields", "", "Alias for -f")
	flag.BoolVar(&args.SkipHeader, "header", false, "Alias for -h")
	flag.BoolVar(&args.Help, "help", false, "Help")
	flag.BoolVar(&args.Version, "version", false, "Version")
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("cutc: ")

	if args.Help {
		fmt.Println(help)
		os.Exit(0)
	}

	if args.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	err := cutc.Run(os.Stdin, os.Stdout, args)
	if err != nil {
		log.Fatalln(err)
	}
}
