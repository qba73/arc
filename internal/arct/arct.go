// Package arct provides functions for parsing arctool log files
// and generating reports.
package arct

=======
package arct

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Marshal ...
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal ...
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// MarshalCSV
var MarshalCSV = func(v interface{}) (io.Reader, error) {
	data :=  
	return bytes.NewReader(), nil

}

var lock sync.Mutex

// Save saves a representation of v to the file in path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := Marshal(v)
	_, err = io.Copy(f, r)

	return err
}

// Load ...
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v)
}

// UnmarshalCSV ...
var UnmarshalCSV = func(r io.Reader, v interface{}) error {
	return nil
}

// LoadReport ...
func LoadReport(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	//data, err := ioutil.ReadFile(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	scaner := bufio.NewScanner(f)

	var lines []string
	var srno, wprn, premiseid string

	for scaner.Scan() {
		l := scaner.Text()
		if strings.HasPrefix(l, "Sr.No") {
			l = strings.ReplaceAll(l, ";", "")
			_, err := fmt.Sscanf(l, "Sr.No = %s WPRN = %s PremiseID = %s", &srno, &wprn, &premiseid)
			if err != nil {
				fmt.Println(err)
			}

			line := fmt.Sprintf("%s,%s,%s\n", srno, wprn, premiseid)
			lines = append(lines, line)
		}
	}

	out := strings.Join(lines, "")
	strings.NewReader(out)

	csv.NewReader(out)

	return nil
}

// Process ...
func Process(data []byte) ([]string, error) {
	return nil, nil
}
