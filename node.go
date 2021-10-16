package dionysus

// node implements Template node
type node struct {
	to        string
	from      string
	bind      string
	staticVal interface{}
	nodes     []*node
	args      Args
}

// Node creates new node instance
func Node() *node {
	return &node{}
}

// Bind sets a binding key. When this node is encoded, by this the bind key,
// a binding object fetched from an encoding interface value
func (e *node) Bind(bind string) *node {
	e.bind = bind

	return e
}

// To sets output node tag name
func (e *node) To(name string) *node {
	e.to = name

	return e
}

// From sets a field name by a binding object from the top-level node
func (e *node) From(from string) *node {
	e.from = from

	return e
}

// StaticVal sets static value for node, priority, if filled
func (e *node) StaticVal(val interface{}) *node {
	e.staticVal = val

	return e
}

// AddNode adds a child node
func (e *node) AddNode(n *node) *node {
	e.nodes = append(e.nodes, n)

	return e
}

// AddAttr adds attribute to node's attributes
func (e *node) AddAttr(a attr) *node {
	e.args = append(e.args, &a)

	return e
}

func (e *node) toMap() map[string]interface{} {
	var args = make([]map[string]interface{}, 0, len(e.args))
	for _, a := range e.args {
		args = append(args, map[string]interface{}{
			"to":        a.to,
			"from":      a.from,
			"staticVal": a.staticVal,
		})
	}

	var nodes = make([]map[string]interface{}, 0, len(e.nodes))
	for _, n := range e.nodes {
		nMap := n.toMap()
		nodes = append(nodes, nMap)
	}

	return map[string]interface{}{
		"to":        e.to,
		"from":      e.from,
		"bind":      e.bind,
		"staticVal": e.staticVal,
		"nodes":     nodes,
		"args":      args,
	}
}

func (e *node) fromMap(m map[string]interface{}) {
	if len(m) == 0 {
		return
	}

	toField := m["to"]
	if toField != nil {
		e.to = toField.(string)
	}

	fromField := m["from"]
	if fromField != nil {
		e.from = fromField.(string)
	}

	bindField := m["bind"]
	if bindField != nil {
		e.bind = bindField.(string)
	}

	staticValField := m["staticVal"]
	e.staticVal = staticValField

	argsField := m["args"]
	if argsField != nil {
		args, ok := argsField.([]interface{})
		if ok {
			e.args = make(Args, 0, len(args))

			for _, iArg := range args {
				a, ok := iArg.(map[string]interface{})
				if !ok {
					continue
				}

				argument := &attr{}

				argTo := a["to"]
				if argTo != nil {
					argument.to = argTo.(string)
				}

				argFrom := a["from"]
				if argFrom != nil {
					argument.from = argFrom.(string)
				}

				staticVal := a["staticVal"]
				argument.staticVal = staticVal

				e.args = append(e.args, argument)
			}
		}
	}

	nodesField := m["nodes"]
	if nodesField == nil {
		return
	}

	nodes, ok := nodesField.([]interface{})
	if !ok {
		return
	}

	e.nodes = make([]*node, 0, len(nodes))
	for _, iNode := range nodes {
		nMap, ok := iNode.(map[string]interface{})
		if !ok {
			continue
		}

		nn := &node{}

		nn.fromMap(nMap)

		e.nodes = append(e.nodes, nn)
	}
}
