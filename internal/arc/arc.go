// Package arc provides functionality for parsing arctool
// log files and generating reports in csv format.
package arc

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// GenerateReport takes a path to the logfile and generate
// csv report in the given output file. If the operation
// is not successfull it returns an error.
func GenerateReport(filein, fileout string) error {
	records, err := loadReportLog(filein)
	if err != nil {
		return err
	}

	fout, err := os.Create(fileout)
	if err != nil {
		return err
	}

	if err := writeCSV(records, fout); err != nil {
		return err
	}
	return nil
}

// loadReportLog knows how to load ArcTool log file.
func loadReportLog(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return processReport(f)
}

// ProcessReport knows how to extract data from a log file
// and return them in a format suitable for writing into csv file.
func processReport(r io.Reader) ([][]string, error) {
	var lines [][]string
	var srno, wprn, premiseid string

	header := []string{"Sr.No", "WPRN", "PremiseID"}
	lines = append(lines, header)

	scaner := bufio.NewScanner(r)
	for scaner.Scan() {
		l := scaner.Text()
		if strings.HasPrefix(l, "Sr.No") {
			l = strings.ReplaceAll(l, ";", "")
			_, err := fmt.Sscanf(l, "Sr.No = %s WPRN = %s PremiseID = %s", &srno, &wprn, &premiseid)
			if err != nil {
				return nil, fmt.Errorf("error when processing log line: %v", err)
			}

			line := []string{srno, wprn, premiseid}
			//line := fmt.Sprintf("%s,%s,%s\n", srno, wprn, premiseid)
			lines = append(lines, line)
		}
	}

	// We don't create a csv report file if the input
	// log data file does not contain data we are interested in.
	if len(lines) == 1 {
		return nil, fmt.Errorf("processed log report does not contain data")
	}

	return lines, nil
}

func writeCSV(records [][]string, w io.Writer) error {
	csvwriter := csv.NewWriter(w)
	if err := csvwriter.WriteAll(records); err != nil {
		return err
	}
	return nil
}
