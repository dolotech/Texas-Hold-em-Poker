/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:40
 * Filename : bench_test.go
 * Description :
 * *******************************************************/
package csv

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"
)

// Benchmark with large CSV data set

func loadData() []byte {
	f, err := os.Open("testdata/ercot-dam.csv.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	return data
}

type Price struct {
	DeliveryDate string
	HourEnding   string
	BusName      string
	LMP          float32
	DSTFlag      bool `true:"Y" false:"N"`
}

func BenchmarkUnmarshal(b *testing.B) {
	b.StopTimer()
	data := loadData()
	b.StartTimer()

	pp := []Price{}

	err := Unmarshal(data, &pp)
	if err != nil {
		panic(err)
	}

	out, err := Marshal(pp)
	if err != nil {
		panic(err)
	}

	if bytes.Equal(data, out) != true {
		panic("wrong results")
	}

	b.Logf("Unmarshal %d rows", len(pp))
	b.Logf("Marshaled %d bytes", len(out))
}
