package arc_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/arc"
)

func TestProcessReportToCSV_ReturnsErrorOnReadError(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToCSV(iotest.ErrReader(errors.New("test error")), io.Discard)
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToJSON_ReturnsErrorOnReadError(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToJSON(iotest.ErrReader(errors.New("test error")), io.Discard)
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToCSV_ReturnsErrorOnInvalidLine(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToCSV(strings.NewReader(invalidData), io.Discard)
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToJSON_ReturnsErrorOnInvalidLine(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToJSON(strings.NewReader(invalidData), io.Discard)
	if err == nil {
		t.Fatal(nil)
	}
}

// errWriter implements writer interface. Calling Write
// method on errWriter will always return error.
type errWriter struct{}

// Write returns writer error. It's used solely for testing.
func (errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("writer error")
}

func TestProcessReportToCSV_ReturnsErrorOnWriteError(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToCSV(strings.NewReader(validData), errWriter{})
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToJSON_ReturnsErrorOnWriteError(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToJSON(strings.NewReader(validData), errWriter{})
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToCSV_ReturnsErrorOnNoData(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToCSV(strings.NewReader(noData), io.Discard)
	if err == nil {
		t.Fatal(err)
	}
}

func TestProcessReportToJSON_ReturnsErrorOnNoData(t *testing.T) {
	t.Parallel()
	err := arc.ProcessReportToJSON(strings.NewReader(noData), io.Discard)
	if err == nil {
		t.Fatal(nil)
	}
}

func TestProcessReportToJSON_ProducesCorrectOutput(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	err := arc.ProcessReportToJSON(strings.NewReader(twoLinesData), buf)
	if err != nil {
		t.Fatal(err)
	}

	want := correctJSONoutput
	got := buf.String()

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestProcessReportToCSV_ProducesCorrectOutput(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	err := arc.ProcessReportToCSV(strings.NewReader(validData), buf)
	if err != nil {
		t.Fatal(err)
	}

	want := correctCSVoutput
	got := buf.String()

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}

}

func TestParseReport_ReadsLogDataAndProducesReport(t *testing.T) {
	t.Parallel()
	want := twoLinesReport
	got, err := arc.ParseReport(strings.NewReader(twoLinesData))
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestParseReport_ReturnsErrorOnNoData(t *testing.T) {
	t.Parallel()
	_, err := arc.ParseReport(strings.NewReader(noData))
	if err == nil {
		t.Fatal(err)
	}
}

func TestParseReport_ReturnsErrorOnInvalidLine(t *testing.T) {
	t.Parallel()
	_, err := arc.ParseReport(strings.NewReader(invalidData))
	if err == nil {
		t.Fatal(err)
	}
}

func TestParser_GeneratesReportInCSVFormat(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}

	p, err := arc.NewParser(
		arc.WithInput(bytes.NewBufferString(validData)),
		arc.WithOutput(buf),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = p.ToCSV()
	if err != nil {
		t.Fatal(err)
	}
	want := correctCSVoutput
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParser_GeneratesReportInJSONFormat(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}

	p, err := arc.NewParser(
		arc.WithInput(bytes.NewBufferString(twoLinesData)),
		arc.WithOutput(buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = p.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	want := correctJSONoutput
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParser_GeneratesReportInCSVFormatWithInputFromArgs(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	args := []string{"testdata/valid_input_data.log"}
	p, err := arc.NewParser(
		arc.WithInputFromArgs(args),
		arc.WithOutput(buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := correctCSVoutput
	err = p.ToCSV()
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParser_GeneratesReportWithInputFromArgsReadingMultipleFiles(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	args := []string{
		"testdata/valid_input_data.log",
		"testdata/valid_input_data2.log",
	}
	p, err := arc.NewParser(
		arc.WithInputFromArgs(args),
		arc.WithOutput(buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := correctCSVoutputFromTwoFiles
	err = p.ToCSV()
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParser_ShowsUsageOnBogusFlag(t *testing.T) {
	t.Parallel()
	args := []string{"-w"}
	_, err := arc.NewParser(
		arc.WithOutput(io.Discard),
		arc.WithInputFromArgs(args),
	)
	if err == nil {
		t.Fatal("want err on bogus flag, got nil")
	}
}

func TestParser_GeneratesReportInCSVFormatOnEmptyInputArgs(t *testing.T) {
	t.Parallel()
	inputBuf := bytes.NewBufferString(validData)
	outputBuf := &bytes.Buffer{}
	p, err := arc.NewParser(
		arc.WithInput(inputBuf),
		arc.WithInputFromArgs([]string{}),
		arc.WithOutput(outputBuf),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := correctCSVoutput
	err = p.ToCSV()
	if err != nil {
		t.Fatal(err)
	}
	got := outputBuf.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// ==============================================================
// Web Service tests

func TestServer_RespondsWithErrorOnNotAllowedMethodInJSONRequest(t *testing.T) {
	t.Parallel()
	handler := arc.NewArcMux()
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/csv", ts.URL))
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusMethodNotAllowed
	got := res.StatusCode

	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestServer_ReturnsOKAndCSVReport(t *testing.T) {
	t.Parallel()
	handler := arc.NewArcMux()

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	part, err := w.CreateFormFile("file", "file.log")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, strings.NewReader(validData))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/csv", b)
	req.Header.Add("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	want := http.StatusOK
	got := res.Code

	if want != got {
		t.Errorf("want %d, got %v", want, got)
	}

	wantBody := correctCSVoutput
	gotBody := res.Body.String()

	if !cmp.Equal(wantBody, gotBody) {
		t.Errorf(cmp.Diff(wantBody, gotBody))
	}
}

func TestServer_ReturnsOKAndJSONReport(t *testing.T) {
	t.Parallel()
	handler := arc.NewArcMux()

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	part, err := w.CreateFormFile("file", "file.log")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, strings.NewReader(twoLinesData))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/json", b)
	req.Header.Add("Content-Type", w.FormDataContentType())

	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	want := http.StatusOK
	got := res.Code

	if want != got {
		t.Errorf("want %d, got %v", want, got)
	}

	wantBody := correctJSONoutput
	gotBody := res.Body.String()

	if !cmp.Equal(wantBody, gotBody) {
		t.Errorf(cmp.Diff(wantBody, gotBody))
	}
}

var (
	invalidData = `—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Sr.No = SKSZPUB0257-1; WPRN = 2607303; PremiseID = 2306982
Sr.No = SKSZPUB0257-2; WPRM = 2607304; PremiseIDx = 3104983
Total number of records to be processed 4
Total number of processed records successfully 4
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
`

	validData = `—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Sr.No = SKSZPUB0257-1; WPRN = 2607303; PremiseID = 2306982
Sr.No = SKSZPUB0257-2; WPRN = 2607304; PremiseID = 3104983
Sr.No = SKSZPUB0257-3; WPRN = 2607305; PremiseID = 5616984
Sr.No = SKSZPUB0257-4; WPRN = 2607306; PremiseID = 1626985
Total number of records to be processed 4
Total number of processed records successfully 4
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
—————————————START————————————19/02/2021 12:53:48
File loading is started for file dfpub0441_DLRZone4_residential_2021-02-18.csv
Sr.No = ALSZPUB0241-1; WPRN = 1507307; PremiseID = 2601986
Sr.No = ALSZPUB0241-2; WPRN = 1507308; PremiseID = 2601987
Sr.No = ALSZPUB0241-3; WPRN = 1507309; PremiseID = 2601988
Total number of records to be processed 3
Total number of processed records successfully 3
Total number of failed records 0
—————————————END—————————————19/02/2021 12:58:21
`

	noData = `—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Total number of records to be processed 0
Total number of processed records successfully 0
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21`

	correctCSVoutput = `Sr.No,WPRN,PremiseID
SKSZPUB0257-1,2607303,2306982
SKSZPUB0257-2,2607304,3104983
SKSZPUB0257-3,2607305,5616984
SKSZPUB0257-4,2607306,1626985
ALSZPUB0241-1,1507307,2601986
ALSZPUB0241-2,1507308,2601987
ALSZPUB0241-3,1507309,2601988
`

	correctCSVoutputFromTwoFiles = `Sr.No,WPRN,PremiseID
SKSZPUB0257-1,2607303,2306982
SKSZPUB0257-2,2607304,3104983
SKSZPUB0257-3,2607305,5616984
SKSZPUB0257-4,2607306,1626985
ALSZPUB0241-1,1507307,2601986
ALSZPUB0241-2,1507308,2601987
ALSZPUB0241-3,1507309,2601988
SKSZPUB0258-1,2607305,4306982
SKSZPUB0258-2,2607306,6104983
SKSZPUB0258-3,2607307,7616984
SKSZPUB0258-4,2607308,5626985
ALSZPUB0243-1,2507307,4601986
ALSZPUB0243-2,2507308,4601987
ALSZPUB0243-3,2507309,4601988
`

	correctJSONoutput = `[{"sr_no":"SKSZPUB0257-1","wprn":"2607303","premise_id":"2306982"},{"sr_no":"SKSZPUB0257-2","wprn":"2607304","premise_id":"3104983"}]`

	twoLinesData = `—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Sr.No = SKSZPUB0257-1; WPRN = 2607303; PremiseID = 2306982
Sr.No = SKSZPUB0257-2; WPRN = 2607304; PremiseID = 3104983
Total number of records to be processed 2
Total number of processed records successfully 2
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
`

	twoLinesReport = arc.Report{
		{SerialNumber: "SKSZPUB0257-1", WPRN: "2607303", PremiseID: "2306982"},
		{SerialNumber: "SKSZPUB0257-2", WPRN: "2607304", PremiseID: "3104983"},
	}
)
