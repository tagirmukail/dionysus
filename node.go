package dionysus

// Node implements Template field
type Node struct {
	to        string
	bind      string
	staticVal interface{}
	nodes     Nodes
	args      Args
}

// NewNode create new Node
func NewNode(
	to string,
	bind string,
	staticVal interface{},
	args Args,
	nodes Nodes,
) *Node {
	return &Node{
		to:        to,
		bind:      bind,
		staticVal: staticVal,
		args:      args,
		nodes:     nodes,
	}
}

// BindTo binds this Node to from data field
// e.g.
// input struct {
//		Foo struct {
//			Bar string `dion:"bar"`
//		} `dion:"foo"`
//	}{
//		Foo: struct {
//			Bar string `dion:"bar"`
//		}{Bar: "a"},
//	}
// e.BindTo("foo.bar")
func (e Node) BindTo(name string) Node {
	e.bind = name

	return e
}

// To initialize binding field name from Node.bind
// e.g.
//
// input
//  struct {
//		Foo struct {
//			Bar string `dion:"bar"`
//		} `dion:"foo"`
//	}{
//		Foo: struct {
//			Bar string `dion:"bar"`
//		}{Bar: "a"},
//	}
//
// e.BindTo("foo.bar")
//
// e.To("test.name")
//
// output
//
// json { "test": { "name": "a" } }
func (e Node) To(name string) Node {
	e.to = name

	return e
}

// StaticVal sets static value for Node, if bind data field is empty
func (e Node) StaticVal(val interface{}) Node {
	e.staticVal = val

	return e
}

// AddNode adds a child Node to the parent Node
func (e Node) AddNode(n Node) Node {
	e.nodes = append(e.nodes, &n)

	return e
}

// AddArg adds argument to arguments for Node
func (e Node) AddArg(a Arg) Node {
	e.args = append(e.args, &a)

	return e
}
