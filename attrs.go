package gotemplconstr

type Attrs []*attr

// attr implement node attribute
type attr struct {
	to        string
	from      string
	staticVal interface{}
}

// Attr creates new struct of attr
func Attr() attr {
	return attr{}
}

// To initialize output attribute name
func (a attr) To(name string) attr {
	a.to = name

	return a
}

// From sets from incoming data value
func (a attr) From(from string) attr {
	a.from = from

	return a
}

// StaticVal sets this static val to attribute, priority if filled
func (a attr) StaticVal(val interface{}) attr {
	a.staticVal = val

	return a
}
