/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:42
 * Filename : example_test.go
 * Description :
 * *******************************************************/
package csv

import (
	"fmt"
	"github.com/golang/glog"
)

type Person struct {
	Name    string  `csv:"Full Name"`
	income  string  // unexported fields are not Unmarshalled
	Age     int     `csv:"-"`      // skip this field
	Address Address `csv:"Street"` // skip this field
}

type Address struct {
	City   string
	Street string
}

func (a *Address) UnmarshalCSV(val string, row *Row) error {
	c, _ := row.Named("City")
	s, _ := row.Named("Street")

	a.Street = s
	a.City = c

	return nil
}

func ExampleUnmarshal() {
	people := []Person{}

	sample := []byte(
		`Full Name,income,Age,City,Street
John Doe,"32,000",45,Brooklyn,"7th Street"
`)

	err := Unmarshal(sample, &people)

	if err != nil {
		glog.Infoln("Error: ", err)
	}

	glog.Infof("%+v", people)

	// Output:
	// [{Name:John Doe income: Age:0 Address:{City:Brooklyn Street:7th Street}}]
}
