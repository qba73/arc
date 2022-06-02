package arc

import (
	"fmt"
	"os"
)

// Values of these vars are passed during the build (Makefile).
var (
	version = ""
	commit  = ""
	date    = ""
)

var usage string = `
arc - ArcTool log processor for generating CSV or JSON data files.

Flags:

-h	"Show help"
-v	"Show version"
-out	"A path to the file you want to generate."

Examples:

	// Generate csv file with data from the log file
	arc < LoaderLogs_19-02-2020.log > report.csv

	// Generate csv file with data from multiple log files
	arc < file1.log file2.log file3.log > report.csv
`

// RunCLI parses arguments and executes program.
func RunCLI() {
	p, err := NewParser(
		WithInputFromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprint(os.Stderr)
		os.Exit(1)
	}
	if p.help {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}
	if p.version {
		fmt.Fprint(os.Stdout, showVersion())
		os.Exit(0)
	}
	p.ToCSV()
}

func showVersion() string {
	return fmt.Sprintf("Version: %s\nGitRef: %s\nBuild Time: %s\n", version, commit, date)
}
