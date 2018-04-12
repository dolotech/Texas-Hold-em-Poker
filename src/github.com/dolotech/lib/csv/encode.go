/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:41
 * Filename : encode.go
 * Description :
 * *******************************************************/
package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Marshaler is an interface for objects which can Marshal themselves into CSV.
type Marshaler interface {
	MarshalCSV() ([]byte, error)
}

type encoder struct {
	*csv.Writer
	buffer *bytes.Buffer
}

// Marshal returns the CSV encoding of i, which must be a slice of struct types.
//
// Marshal traverses the slice and encodes the primative values.
//
// The first row of the CSV output is a header row. The column names are based
// on the field name.  If a different name is required a struct tag can be used
// to define a new name.
//
//   Field string `csv:"Column Name"`
//
// To skip encoding a field use the "-" as the tag value.
//
//   Field string `csv:"-"`
//
// Boolean fields can use string values to define true or false.
//   Bool bool `true:"Yes" false:"No"`
func Marshal(i interface{}) ([]byte, error) {
	// validate the interface
	// create a new encoder
	//   assing the cfields
	// get the headers
	// encoder each row

	data := reflect.ValueOf(i)

	if data.Kind() != reflect.Slice {
		return []byte{}, errors.New("only slices can be marshalled")
	}

	el := data.Index(0)
	enc, err := newEncoder(el)

	if err != nil {
		return []byte{}, err
	}

	err = enc.encodeAll(data)

	if err != nil {
		return []byte{}, err
	}

	enc.Flush()
	return enc.buffer.Bytes(), nil
}

func newEncoder(el reflect.Value) (*encoder, error) {
	b := bytes.NewBuffer([]byte{})

	enc := &encoder{
		buffer: b,
		Writer: csv.NewWriter(b),
	}

	err := enc.Write(colNames(el.Type()))

	return enc, err
}

// colNames takes a struct and returns the computed columns names for each
// field.
func colNames(t reflect.Type) (out []string) {
	l := t.NumField()

	for x := 0; x < l; x++ {
		f := t.Field(x)
		h, ok := fieldHeaderName(f)
		if ok {
			out = append(out, h)
		}
	}

	return
}

// encodeAll iterates over each item in data, encoder it then writes it
func (enc *encoder) encodeAll(data reflect.Value) error {
	n := data.Len()
	for c := 0; c < n; c++ {
		row, err := enc.encodeRow(data.Index(c))

		if err != nil {
			return err
		}

		err = enc.Write(row)

		if err != nil {
			return err
		}
	}

	return nil
}

// encodes a struct into a CSV row
func (enc *encoder) encodeRow(v reflect.Value) ([]string, error) {

	var row []string
	// TODO env.columns should map to a cfield
	// iterate over each cfield and encode with it
	l := v.Type().NumField()

	for x := 0; x < l; x++ {
		fv := v.Field(x)
		st := v.Type().Field(x).Tag

		if st.Get("csv") == "-" {
			continue
		}
		o := enc.encodeCol(fv, st)
		row = append(row, o)
	}

	return row, nil
}

// Returns the string representation of the field value
func (enc *encoder) encodeCol(fv reflect.Value, st reflect.StructTag) string {
	switch fv.Kind() {
	case reflect.String:
		return fv.String()
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return fmt.Sprintf("%v", fv.Int())
	case reflect.Float32:
		return encodeFloat(32, fv)
	case reflect.Float64:
		return encodeFloat(64, fv)
	case reflect.Bool:
		return encodeBool(fv.Bool(), st)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%v", fv.Uint())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%+.3g", fv.Complex())
	case reflect.Interface:
		return encodeInterface(fv, st)
	case reflect.Struct:
		return encodeInterface(fv, st)
	default:
		panic(fmt.Sprintf("Unsupported type %s", fv.Kind()))
	}

	return ""
}

func encodeFloat(bits int, f reflect.Value) string {
	return strconv.FormatFloat(f.Float(), 'g', -1, bits)
}

func encodeBool(b bool, st reflect.StructTag) string {
	v := strconv.FormatBool(b)
	tv := st.Get(v)

	if tv != "" {
		return tv
	}
	return v
}

func encodeInterface(fv reflect.Value, st reflect.StructTag) string {
	marshalerType := reflect.TypeOf(new(Marshaler)).Elem()

	if fv.Type().Implements(marshalerType) {
		m := fv.Interface().(Marshaler)
		b, err := m.MarshalCSV()
		if err != nil {
			return ""
		}
		return string(b)
	}

	return ""
}
