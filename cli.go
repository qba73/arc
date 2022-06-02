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

var (
	usageCSV string = `
arc2csv - ArcTool log processor for parsing log files to CSV format.

Flags:

-h	"Show help"
-v	"Show version"

Examples:

	// Generate csv file with data from the log file
	arc2csv < LoaderLogs_19-02-2020.log > report.csv

	// Generate csv file with data from multiple log files
	arc2csv < file1.log file2.log file3.log > report.csv
`

	usageJSON string = `
arc2json - ArcTool log processor for parsing log files to JSON format.

Flags:

-h	"Show help"
-v	"Show version"

Examples:

	// Generate csv file with data from the log file
	arc2json < LoaderLogs_19-02-2020.log > report.json

	// Generate csv file with data from multiple log files
	arc2json < file1.log file2.log file3.log > report.json
`
)

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
		fmt.Fprint(os.Stdout, usageCSV)
		os.Exit(0)
	}
	if p.version {
		fmt.Fprint(os.Stdout, showVersion())
		os.Exit(0)
	}
	p.ToCSV()
}

func RunJSONCLI() {
	p, err := NewParser(
		WithInputFromArgs(os.Args[1:]),
	)
	if err != nil {
		fmt.Fprint(os.Stderr)
		os.Exit(1)
	}
	if p.help {
		fmt.Fprint(os.Stdout, usageJSON)
		os.Exit(0)
	}
	if p.version {
		fmt.Fprint(os.Stdout, showVersion())
		os.Exit(0)
	}
	p.ToJSON()
}

func showVersion() string {
	return fmt.Sprintf("Version: %s\nGitRef: %s\nBuild Time: %s\n", version, commit, date)
}
