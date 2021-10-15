package dionysus

type Args []*arg

// arg implement node argument
type arg struct {
	to        string
	from      string
	staticVal interface{}
}

func Arg() arg {
	return arg{}
}

// NewArg creates new instance of arg
func NewArg(to string, staticVal interface{}) *arg {
	return &arg{
		to:        to,
		staticVal: staticVal,
	}
}

// To initialize out binding field name from arg.bind
func (a arg) To(name string) arg {
	a.to = name

	return a
}

func (a arg) From(from string) arg {
	a.from = from

	return a
}

// StaticVal sets this static val to argument if bind data field is empty
func (a arg) StaticVal(val interface{}) arg {
	a.staticVal = val

	return a
}
