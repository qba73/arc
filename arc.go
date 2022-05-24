// Package arc provides functionality for parsing arctool
// log files and generating reports in the csv format.
package arc

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// GenerateReport takes a path to the logfile and generate
// csv report in the given output file. If the operation
// is not successful it returns an error.
func GenerateReport(filein, fileout string) error {
	fin, err := os.Open(filein)
	if err != nil {
		return fmt.Errorf("opening log file: %s, err: %v", filein, err)
	}
	defer fin.Close()

	fout, err := os.Create(fileout)
	if err != nil {
		return fmt.Errorf("creating output file: %s, err: %v", fileout, err)
	}
	defer fout.Close()

	return ProcessReportToCSV(fin, fout)
}

// isLineWithData holds logic to verify if
// the string holds the data for processing.
func isLineWithData(l string) bool {
	return strings.HasPrefix(l, "Sr.No")
}

type reportLine struct {
	SerialNumber string `json:"sr_no"`
	WPRN         string `json:"wprn"`
	PremiseID    string `json:"premise_id"`
}

// Report represents data ready to save in JSON or CSV formats.
type Report []reportLine

// CSV formats report data in the CSV format.
func (r Report) CSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	rep := [][]string{{"Sr.No", "WPRN", "PremiseID"}}
	for _, l := range r {
		rep = append(rep, []string{l.SerialNumber, l.WPRN, l.PremiseID})
	}
	return writer.WriteAll(rep)
}

// JSON formats report data in the JSON format.
func (r Report) JSON(w io.Writer) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// ParseReport takes a reader and returs a report. The report is a
// slice of structs containing parsed data from the reader.
// ParseReport willreturn error on malformed input data
// or when the reader does not contain data to parse.
func ParseReport(r io.Reader) (Report, error) {
	var report Report
	var srno, wprn, premiseid string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()

		if !isLineWithData(l) {
			continue
		}

		// We need to clean ';' to prepare string for Sscanf function.
		l = strings.ReplaceAll(l, ";", "")
		_, err := fmt.Sscanf(l, "Sr.No = %s WPRN = %s PremiseID = %s", &srno, &wprn, &premiseid)
		if err != nil {
			return nil, fmt.Errorf("processing log line: %s, %v", l, err)
		}

		line := reportLine{
			SerialNumber: srno,
			WPRN:         wprn,
			PremiseID:    premiseid,
		}
		report = append(report, line)
	}

	if len(report) == 0 {
		return nil, errors.New("no data in the input file")
	}
	return report, nil
}

// ProcessReportToJSON takes a reader and writes report in JSON
// format to the writer. If parsing is not successful it returns an error.
func ProcessReportToJSON(r io.Reader, w io.Writer) error {
	data, err := ParseReport(r)
	if err != nil {
		return err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// ProcessReportToCSV takes a reader and writes report in CSV format
// to the writer. If parsing is not successfull it returns an error.
func ProcessReportToCSV(r io.Reader, w io.Writer) error {
	records, err := ParseReport(r)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(w)
	rep := [][]string{{"Sr.No", "WPRN", "PremiseID"}}
	for _, l := range records {
		rep = append(rep, []string{l.SerialNumber, l.WPRN, l.PremiseID})
	}
	return writer.WriteAll(rep)
}

// RunCLI parses arguments and executes program.
func RunCLI() {
	fmt.Println("running cli...")
}

// uploadFile is the handler responsible for processing
// raw data files. It returns
func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading report file for processing\n")
	w.Header().Add("Content-Type", "application/json")

	var maxSize int64 = 1024 * 1024

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(w, "report file is too big", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = ProcessReportToJSON(file, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RunServer takes config params and runs arc web werver.
func RunServer() {
	fmt.Println("server starting")
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
