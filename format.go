package atlas_sdk

import "strings"

const (
	sqlFmt  = "{{ sql . }}"
	jsonFmt = "{{ json . }}"
)

// Format
//
// Use SQLFormat or JSONFormat, and you can create your own by NewFormat.
// String must correspond to text/template.
type Format interface {
	GoFormat() string
}

// NewFormat creates a new instance of the customFormat struct
// which implements the Format interface. The fmt parameter is
// used to set the format of the custom format.
func NewFormat(fmt string) Format {
	return &customFormat{fmt: strings.Trim(fmt, `"`)}
}

// customFormat is a type that represents a custom format used in the Format interface.
// It stores the format string in the fmt field.
type customFormat struct {
	fmt string
}

// GoFormat returns the formatted string of a customFormat object.
func (f *customFormat) GoFormat() string {
	return f.fmt
}

type SQLFormat struct{}

// GoFormat returns the SQL format template used by the SQLFormat struct.
func (f *SQLFormat) GoFormat() string {
	return sqlFmt
}

// JSONFormat represents a type that formats data into JSON format using the text/template package.
type JSONFormat struct{}

// GoFormat returns the Go format for JSON.
func (f *JSONFormat) GoFormat() string {
	return jsonFmt
}
