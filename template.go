package gotemplconstr

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
)

// FileType is output file format
type FileType uint16

const (
	JSON FileType = iota
	XML
	CSV
	TSV
	YAML
)

func (f FileType) String() string {
	switch f {
	case JSON:
		return "json"
	case XML:
		return "xml"
	case CSV:
		return "csv"
	case TSV:
		return "tsv"
	case YAML:
		return "yaml"
	default:
		return strconv.FormatUint(uint64(f), 10)
	}
}

// Template implements template
type Template struct {
	outputType FileType
	node       *node

	printer printer
}

// NewTemplate creates new instance of Template
func NewTemplate() *Template {
	return &Template{}
}

// AddNode adds a general node to the Template
func (t *Template) AddNode(n *node) *Template {
	t.node = n

	return t
}

// ToOutputFileType sets output file type for this Template
func (t *Template) ToOutputFileType(ft FileType) *Template {
	t.outputType = ft

	return t
}

func (t *Template) FileType() FileType {
	return t.outputType
}

func (t *Template) newPrinter(w io.Writer) {
	buff := bytes.Buffer{}

	t.printer = printer{
		Writer: w,
		tmp:    &buff,
	}
}

func (t *Template) NewEncoder(w io.Writer) *Template {
	t.newPrinter(w)

	return t
}

// MarshalJSON is represents template in json format
func (t *Template) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(map[string]interface{}{
		"outputType": t.outputType,
		"node":       t.node.toMap(),
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalJSON is restores template from json
func (t *Template) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})

	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	ot := m["outputType"]
	if ot != nil {
		oType, ok := ot.(float64)
		if ok {
			t.outputType = FileType(oType)
		}
	}

	n := m["node"]
	if n != nil {
		nd, ok := n.(map[string]interface{})
		if ok {
			t.node = &node{}
			t.node.fromMap(nd)
		}
	}

	return nil
}
