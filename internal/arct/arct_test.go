package arct_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/arct/internal/arct"
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

type line struct {
	srno      string
	wprn      string
	premiseid string
}

func TestLoadFile(t *testing.T) {

	/*
		items := [][]string{
			{"Sr.No", "WPRN", "PremiseID"},
		}
	*/

	testFile := "testdata/BulkLoaderLogs.txt"
	got, err := arct.LoadReport(testFile)

	if err != nil {
		t.Fatalf("error loading file: %s", testFile)
	}

	if !cmp.Equal(string(got), content) {
		t.Errorf(cmp.Diff(string(got), content))
	}
}

func TestProcess(t *testing.T) {

	data, err := arct.LoadReport("testdata/BulkLoaderLogs.txt")
	if err != nil {
		t.Fatalf()
	}

}

/*
	scaner := bufio.NewScanner(strings.NewReader(string(data)))

	for scaner.Scan() {
		//var s1, s2, s3 string
		l := scaner.Text()
		if strings.HasPrefix(l, "Sr.No") {

			var srno, wprn, premiseid string
			l := strings.ReplaceAll(l, ";", "")
			n, err := fmt.Sscanf(l, "Sr.No = %s WPRN = %s PremiseID = %s", &srno, &wprn, &premiseid)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(n)
			itm := []string{srno, wprn, premiseid}

			items = append(items, itm)
		}
	}

	f, err := os.Create("out.csv")
	if err != nil {
		log.Fatalln(err)
	}

	w := csv.NewWriter(f)

	for _, record := range items {
		if err := w.Write(record); err != nil {
			log.Fatalln(err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatalln(err)
	}

*/
