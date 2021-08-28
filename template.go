package Dionysus

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
		return ""
	}
}

// Template implements output file template
type Template struct {
	outputType FileType
	node       *Node
}

// NewTemplate creates new instance of Template
func NewTemplate(outputType FileType, n *Node) *Template {
	return &Template{
		outputType: outputType,
		node:       n,
	}
}

// AddNode adds a child Node to the Template
func (t *Template) AddNode(n Node) *Node {
	t.node = &n

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
