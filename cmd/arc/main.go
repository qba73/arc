package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qba73/arct/internal/arc"
)

func main() {
	var filein, fileout string
	flag.StringVar(&filein, "input", "", "ArcTool log file to process")
	flag.StringVar(&fileout, "output", "", "Output CSV file")
	flag.Parse()
	verifyFlags(filein, fileout)

	if err := run(filein, fileout); err != nil {
		fmt.Println("Encountered an error, exiting")
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
