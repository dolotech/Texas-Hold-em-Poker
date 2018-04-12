/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-07-07 23:41
 * Filename : csv.go
 * Description :
 * *******************************************************/
// Package csv provides `Marshal` and `UnMarshal` encoding functions for CSV(Comma Seperated Value) data.
// This package is built on the the standard library's encoding/csv.
package csv

import (
	"reflect"
)

// fieldHeaderName returns the header name to use for the given StructField
// This can be a user defined name (via the Tag) or a default name.
func fieldHeaderName(f reflect.StructField) (string, bool) {
	h := f.Tag.Get("csv")

	if h == "-" {
		return "", false
	}

	// If there is no tag set, use a default name
	if h == "" {
		return f.Name, true
	}

	return h, true
}
