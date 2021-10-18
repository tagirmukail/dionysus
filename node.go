package gotemplconstr

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
func (nn *node) Bind(bind string) *node {
	nn.bind = bind

	return nn
}

// To sets output node tag name
func (nn *node) To(name string) *node {
	nn.to = name

	return nn
}

// From sets a field name by a binding object from the top-level node
func (nn *node) From(from string) *node {
	nn.from = from

	return nn
}

// StaticVal sets static value for node, priority, if filled
func (nn *node) StaticVal(val interface{}) *node {
	nn.staticVal = val

	return nn
}

// AddNode adds a child node
func (nn *node) AddNode(ns ...*node) *node {
	nn.nodes = append(nn.nodes, ns...)

	return nn
}

// AddAttr adds attribute to node's attributes
func (nn *node) AddAttr(a attr) *node {
	nn.args = append(nn.args, &a)

	return nn
}

func (nn *node) toMap() map[string]interface{} {
	var args = make([]map[string]interface{}, 0, len(nn.args))
	for _, a := range nn.args {
		args = append(args, map[string]interface{}{
			"to":        a.to,
			"from":      a.from,
			"staticVal": a.staticVal,
		})
	}

	var nodes = make([]map[string]interface{}, 0, len(nn.nodes))
	for _, n := range nn.nodes {
		nMap := n.toMap()
		nodes = append(nodes, nMap)
	}

	return map[string]interface{}{
		"to":        nn.to,
		"from":      nn.from,
		"bind":      nn.bind,
		"staticVal": nn.staticVal,
		"nodes":     nodes,
		"args":      args,
	}
}

func (nn *node) fromMap(m map[string]interface{}) {
	if len(m) == 0 {
		return
	}

	toField := m["to"]
	if toField != nil {
		nn.to = toField.(string)
	}

	fromField := m["from"]
	if fromField != nil {
		nn.from = fromField.(string)
	}

	bindField := m["bind"]
	if bindField != nil {
		nn.bind = bindField.(string)
	}

	staticValField := m["staticVal"]
	nn.staticVal = staticValField

	argsField := m["args"]
	if argsField != nil {
		args, ok := argsField.([]interface{})
		if ok {
			nn.args = make(Args, 0, len(args))

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

				nn.args = append(nn.args, argument)
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

	nn.nodes = make([]*node, 0, len(nodes))
	for _, iNode := range nodes {
		nMap, ok := iNode.(map[string]interface{})
		if !ok {
			continue
		}

		n := &node{}

		n.fromMap(nMap)

		nn.nodes = append(nn.nodes, n)
	}
}
