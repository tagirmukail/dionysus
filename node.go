package gotemplconstr

// node implements Template node
type node struct {
	to        string
	from      string
	bind      string
	staticVal interface{}
	nodes     []*node
	attrs     Attrs
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
	nn.attrs = append(nn.attrs, &a)

	return nn
}

func (nn *node) toMap() map[string]interface{} {
	var attrs = make([]map[string]interface{}, 0, len(nn.attrs))
	for _, a := range nn.attrs {
		attrs = append(attrs, map[string]interface{}{
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
		"attrs":     attrs,
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

	attrsField := m["attrs"]
	if attrsField != nil {
		attrs, ok := attrsField.([]interface{})
		if ok {
			nn.attrs = make(Attrs, 0, len(attrs))

			for _, iAttr := range attrs {
				a, ok := iAttr.(map[string]interface{})
				if !ok {
					continue
				}

				attribute := &attr{}

				attrTo := a["to"]
				if attrTo != nil {
					attribute.to = attrTo.(string)
				}

				attrFrom := a["from"]
				if attrFrom != nil {
					attribute.from = attrFrom.(string)
				}

				staticVal := a["staticVal"]
				attribute.staticVal = staticVal

				nn.attrs = append(nn.attrs, attribute)
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
