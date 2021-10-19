package gotemplconstr

import (
	"encoding/xml"
	"reflect"
	"strings"
	"time"
)

var replaceXMLSymbols = []replaceSymbol{
	{
		from: `"`,
		to:   "&quot;",
	},
	{
		from: "&",
		to:   "&amp;",
	},
	{
		from: ">",
		to:   "&gt;",
	},
	{
		from: "<",
		to:   "&lt;",
	},
	{
		from: "'",
		to:   "&apos;",
	},
}

type replaceSymbol struct {
	from string
	to   string
}

func xmlReplaceSymbols(s string) string {
	for _, repS := range replaceXMLSymbols {
		s = strings.ReplaceAll(s, repS.from, repS.to)
	}

	return s
}

func (t *Template) encodeXML(globVal reflect.Value) error {
	_, err := t.printer.WriteString(xml.Header)
	if err != nil {
		return err
	}

	err = t.node.encodeXML(&t.printer, globVal, reflect.Value{}, false)
	if err != nil {
		return err
	}

	return t.printer.flush(true)
}

func (nn *node) encodeXML(p *printer, globVal, inVal reflect.Value, isItem bool) error {
	switch {
	case nn == nil:
		return nil
	case nn.staticVal != nil:
		val := reflect.ValueOf(nn.staticVal)
		err := nn.startXML(p, val)
		if err != nil {
			return err
		}

		kind := val.Kind()
		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
		}

		if kind == reflect.Map {
			return ErrStaticValOnlySimpleType
		}

		err = p.xmlWriteReflectVal(val)
		if err != nil {
			return err
		}

		return nn.endXML(p)
	case nn.bind != "":
		val, ok := getVal(nn.bind, globVal)
		if !ok {
			return nil
		}

		kind := val.Kind()
		tp := val.Type()

		err := nn.startXML(p, val)
		if err != nil {
			return err
		}

		if (kind == reflect.Slice || kind == reflect.Array) && tp.Elem().Kind() != reflect.Uint8 {
			if len(nn.nodes) == 0 {
				return nn.endXML(p)
			}

			nod := nn.nodes[0]

			for i, n := 0, val.Len(); i < n; i++ {
				vv := val.Index(i)

				err = nod.encodeXML(p, globVal, vv, true)
				if err != nil {
					return err
				}
			}

			return nn.endXML(p)
		}

		if kind == reflect.Ptr {
			val = val.Elem()
		}

		if kind == reflect.Struct {
			_, ok := val.Interface().(time.Time)
			if ok {
				return ErrBindCantTime
			}

			for _, n := range nn.nodes {
				err = n.encodeXML(p, globVal, val, false)
				if err != nil {
					return err
				}
			}
		}

		return nn.endXML(p)
	case isItem && inVal.IsValid() && !inVal.IsZero():
		err := nn.startXML(p, inVal)
		if err != nil {
			return err
		}

		if len(nn.nodes) > 0 {
			for _, n := range nn.nodes {
				err = n.encodeXML(p, globVal, inVal, false)
				if err != nil {
					return err
				}
			}
		} else {
			fieldVal := inVal.FieldByName(nn.from)
			if fieldVal.Kind() == reflect.Ptr {
				fieldVal = fieldVal.Elem()
			}

			if fieldVal.Kind() == reflect.Map {
				return ErrValOnlySimpleType
			}

			err = p.xmlWriteReflectVal(fieldVal)
			if err != nil {
				return err
			}
		}

		return nn.endXML(p)
	case nn.from != "" && inVal.IsValid() && !inVal.IsZero():
		err := nn.startXML(p, inVal)
		if err != nil {
			return err
		}

		fieldVal := inVal.FieldByName(nn.from)
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		if fieldVal.Kind() == reflect.Map {
			return ErrValOnlySimpleType
		}

		err = p.xmlWriteReflectVal(fieldVal)
		if err != nil {
			return err
		}

		return nn.endXML(p)
	default:
		val := reflect.Value{}

		err := nn.startXML(p, val)
		if err != nil {
			return err
		}

		for _, n := range nn.nodes {
			err = n.encodeXML(p, globVal, val, false)
			if err != nil {
				return err
			}

		}

		return nn.endXML(p)
	}
}

func (nn *node) startXML(p *printer, val reflect.Value) error {
	if strings.TrimSpace(nn.to) == "" {
		return ErrToFiledEmpty
	}

	_, err := p.WriteRune('<')
	if err != nil {
		return err
	}

	_, err = p.WriteString(nn.to)
	if err != nil {
		return err
	}

	for _, a := range nn.attrs {
		err = a.encodeXML(p, val)
		if err != nil {
			return err
		}
	}

	_, err = p.WriteRune('>')
	if err != nil {
		return err
	}

	return nil
}

func (nn *node) endXML(p *printer) error {
	if strings.TrimSpace(nn.to) == "" {
		return ErrToFiledEmpty
	}

	_, err := p.WriteString("</")
	if err != nil {
		return err
	}

	_, err = p.WriteString(nn.to)
	if err != nil {
		return err
	}

	_, err = p.WriteRune('>')
	if err != nil {
		return err
	}

	return nil
}

func (a *attr) encodeXML(p *printer, val reflect.Value) error {
	_, err := p.WriteRune(' ')
	if err != nil {
		return err
	}

	_, err = p.WriteString(a.to)
	if err != nil {
		return err
	}

	_, err = p.WriteString(`="`)
	if err != nil {
		return err
	}

	if a.staticVal != nil {
		statVal := reflect.ValueOf(a.staticVal)
		kind := statVal.Kind()
		if kind == reflect.Ptr {
			statVal = statVal.Elem()
			kind = statVal.Kind()
		}

		if kind == reflect.Map {
			return ErrStaticValOnlySimpleType
		}

		err = p.xmlWriteReflectVal(statVal)
		if err != nil {
			return err
		}

		_, err = p.WriteString(`"`)
		if err != nil {
			return err
		}

		return nil
	}

	kind := val.Kind()

	if kind == reflect.Ptr {
		val = val.Elem()
		kind = val.Kind()
	}

	var fieldVal reflect.Value
	if kind == reflect.Struct {
		fieldVal = val.FieldByName(a.from)
		if !fieldVal.IsValid() {
			return ErrInvalidField
		}
	} else {
		fieldVal = val
	}

	s, b, err := marshalSimpleVal(fieldVal)
	if err != nil {
		return err
	} else if b != nil {
		_, err = p.Write(b)
		if err != nil {
			return err
		}
	} else {
		_, err = p.WriteString(xmlReplaceSymbols(s))
		if err != nil {
			return err
		}
	}

	_, err = p.WriteString(`"`)
	if err != nil {
		return err
	}

	return nil
}

func (p printer) xmlWriteReflectVal(val reflect.Value) (err error) {
	kind := val.Kind()

	if kind == reflect.Struct {
		t, ok := val.Interface().(time.Time)
		if !ok {
			return ErrStaticValOnlySimpleType
		}

		_, err = p.WriteString(t.String())
		if err != nil {
			return err
		}

		return nil
	}

	s, b, err := marshalSimpleVal(val)
	if err != nil {
		return err
	} else if b != nil {
		_, err = p.Write(b)
		if err != nil {
			return err
		}
	} else {
		_, err = p.WriteString(xmlReplaceSymbols(s))
		if err != nil {
			return err
		}
	}

	return nil
}
