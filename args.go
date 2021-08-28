package Dionysus

type Args []*Arg

// Arg implement Node argument
type Arg struct {
	to        string
	bind      string
	staticVal interface{}
}

// NewArg creates new instance of Arg
func NewArg(to, bind string, staticVal interface{}) *Arg {
	return &Arg{
		to:        to,
		bind:      bind,
		staticVal: staticVal,
	}
}

// BindTo binds this Arg to from data field
func (a Arg) BindTo(name string) Arg {

	a.bind = name

	return a
}

// To initialize out binding field name from Arg.bind
func (a Arg) To(name string) Arg {
	a.to = name

	return a
}

// StaticVal sets this static val to argument if bind data field is empty
func (a Arg) StaticVal(val interface{}) Arg {
	a.staticVal = val

	return a
}
