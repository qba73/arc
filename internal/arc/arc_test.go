package arc

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var content = `—————————————START————————————19/02/2021 12:51:06
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
—————————————END—————————————19/02/2021 12:58:21`

var contentNoRecords = `—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Total number of records to be processed 0
Total number of processed records successfully 0
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
—————————————START————————————19/02/2021 12:53:48
File loading is started for file dfpub0441_DLRZone4_residential_2021-02-18.csv
Total number of records to be processed 0
Total number of processed records successfully 0
Total number of failed records 0
—————————————END—————————————19/02/2021 12:58:21`

func TestProcessReport(t *testing.T) {
	data := strings.NewReader(content)

	got, err := processReport(data)

	t.Run("Number of records", func(t *testing.T) {
		wantError := false
		if (err != nil) != wantError {
			t.Fatalf("error processing report: %s", err)
		}

		wantLen := 8
		if len(got) != wantLen {
			t.Errorf("ProcessReport() got %d lines in report, want: %d", len(got), wantLen)
		}
	})

	t.Run("Header", func(t *testing.T) {
		wantError := false
		if (err != nil) != wantError {
			t.Fatalf("error processing report: %s", err)
		}

		wantHeader := []string{"Sr.No", "WPRN", "PremiseID"}
		if !cmp.Equal(got[0], wantHeader) {
			t.Errorf("%s", cmp.Diff(got[0], wantHeader))
		}
	})

	t.Run("Record", func(t *testing.T) {
		wantRecord := []string{"SKSZPUB0257-1", "2607303", "2306982"}
		if !cmp.Equal(got[1], wantRecord) {
			t.Errorf("%s", cmp.Diff(got[1], wantRecord))
		}
	})
}
