package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qba73/arct/internal/arc"
)

// Variables used during the go build.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var filein, fileout string
	var ver bool
	flag.StringVar(&filein, "in", "", "ArcTool log file to process")
	flag.StringVar(&fileout, "out", "", "Output CSV file")
	flag.BoolVar(&ver, "version", false, "Show version")
	flag.Parse()

	flag.Usage = func() {
		fmt.Println(`
Arct - A simple ArcTool log processor for generating CSV data files.

Flags:

-version	"Show arct version and exit"
-in	"A path to the log file generated by the ArcTool."
-out	"A path to the csv file you want to generate."

Usage:

	// To generate csv file with data from the log file
	arct -in=LoaderLogs_19-02-2020.txt -out=data.csv`)
	}

	if ver {
		showVersion()
		os.Exit(0)
	}

	verifyFlags(filein, fileout)

	if err := run(filein, fileout); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(fin, fout string) error {
	if err := arc.GenerateReport(fin, fout); err != nil {
		return err
	}
	return nil
}

func verifyFlags(in, out string) {
	if in == "" || out == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("arct, version: %s\nbuild tag: %s\ndate: %s", version, commit, date)
}
