package dionysus

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
	YML
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
	case YML:
		return "yml"
	default:
		return strconv.FormatUint(uint64(f), 10)
	}
}

// Template implements output file template
type Template struct {
	outputType FileType
	node       *node

	printer printer
}

// NewTemplate creates new instance of Template
func NewTemplate(outputType FileType, n *node) *Template {
	return &Template{
		outputType: outputType,
		node:       n,
	}
}

// AddNode adds a child node to the Template
func (t *Template) AddNode(n *node) *node {
	t.node = n

	return t.node
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
