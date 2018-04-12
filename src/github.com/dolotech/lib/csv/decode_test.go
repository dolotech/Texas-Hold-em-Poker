/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:41
 * Filename : decode_test.go
 * Description :
 * *******************************************************/
package csv

import (
	"reflect"
	"testing"
)

type Q struct {
	String    string
	Int       int
	unxported int
	Bool      bool `true:"Yes" false:"No"`
	Float32   float32
	Float64   float64
	Complex64 complex64 `csv:"C64"`
}

func TestUnmarshal(t *testing.T) {
	doc := []byte(`String,Int,unexported,Bool,Float32,Float64,C64
John,23,1,Yes,32.2,64.1,1
Jane,27,2,No,33.1,65.1,2
Bill,28,3,Yes,34.7,65.1,3`)

	pp := []Q{}

	Unmarshal(doc, &pp)

	if len(pp) != 3 {
		t.Errorf("Incorrect record length: %d", len(pp))
	}

	assert := func(e, a interface{}) {
		if e != a {
			t.Errorf("expected (%s) got (%s)", e, a)
		}
	}

	strs := []string{"John", "Jane", "Bill"}
	ints := []int{23, 27, 28}
	bools := []bool{true, false, true}
	f32s := []float32{32.2, 33.1, 34.7}
	f64s := []float64{64.1, 65.1, 65.1}

	for i, p := range pp {
		assert(strs[i], p.String)
		assert(ints[i], p.Int)
		assert(bools[i], p.Bool)
		assert(f32s[i], p.Float32)
		assert(f64s[i], p.Float64)
	}

}

func TestMarshalErrors(t *testing.T) {
	doc := []byte(`Name,Age`)
	err := Unmarshal(doc, []Q{})
	if err == nil {
		t.Error("No error generated for non-pointer")
	}

	err = Unmarshal(doc, &Q{})
	if err == nil {
		t.Error("No error generated for non-slice")
	}

	pp := []Q{}
	err = Unmarshal(doc, &pp)
	if err != nil {
		t.Error("Error returned when not expected:", err)
	}
}

func TestExportedFields(t *testing.T) {
	type s struct {
		Name string
		Age  int `csv:"Age"`
		priv int `csv:"-"`
	}

	fs := exportedFields(reflect.TypeOf(s{}))

	if len(fs) != 2 {
		t.Error("Incorrect number of exported fields 2 expected got %d", len(fs))
	}

	if fs[0].Name != "Name" || fs[1].Name != "Age" {
		t.Error("Incorrect returned fields")
	}
}

type T struct {
	Name    string
	age     string // unexported, should not be included
	Addr    string `csv:"Address"`
	NoMatch int    // public, but no match in the CSV headers
}

func TestMapFields(t *testing.T) {
	dec := &decoder{Type: reflect.TypeOf(T{})}
	cols := []string{
		"Name",
		"age", // should not match since the 'age' field is not exported
		"Address",
	}

	dec.mapFieldsToCols(cols)

	fm := dec.cfields
	if len(fm) != 2 {
		t.Errorf("Expected length of 2, got %d", len(fm))
	}

	for i, n := range []int{0, 2} {
		if fm[i].colIndex != n {
			t.Errorf("expected colIndex of %d got %d", fm[i].colIndex, n)
		}
	}
}

type Um struct {
	V string
}

type U struct {
	Name Um
}

func (u *Um) UnmarshalCSV(val string, row *Row) error {
	v, _ := row.Named("Age")
	u.V = row.At(0) + " " + v
	return nil
}

func TestCustomUnMarshaller(t *testing.T) {
	doc := `Name,Age
Jay,23`

	oo := []U{}

	Unmarshal([]byte(doc), &oo)

	if oo[0].Name.V != "Jay 23" {
		t.Errorf("custom unmarshal did not work (%s)", oo[0].Name.V)
	}
}
