package dionysus

import (
	"encoding/xml"
	"errors"
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

func (t *Template) encodeXML(globVal reflect.Value) error {
	_, err := t.printer.WriteString(xml.Header)
	if err != nil {
		return err
	}

	err = t.node.encodeXML(&t.printer, globVal, reflect.Zero(reflect.TypeOf(0)), false)
	if err != nil {
		return err
	}

	return t.printer.flush(true)
}

func (e *node) encodeXML(p *printer, globVal, inVal reflect.Value, isItem bool) error {
	switch {
	case e.staticVal != nil:
		val := reflect.ValueOf(e.staticVal)
		err := e.startXML(p, val)
		if err != nil {
			return err
		}

		kind := val.Kind()
		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
		}

		if kind == reflect.Map {
			return errors.New("static val must be only simple type")
		}

		if kind == reflect.Struct {
			iVal := val.Interface()

			t, ok := iVal.(time.Time)
			if !ok {
				return errors.New("static val must be only simple type")
			}

			_, err = p.WriteString(t.String())
			if err != nil {
				return err
			}
		} else {
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
		}

		return e.endXML(p)
	case e.bind != "":
		val, ok := getVal(e.bind, globVal)
		if !ok {
			return nil
		}

		kind := val.Kind()
		tp := val.Type()

		err := e.startXML(p, val)
		if err != nil {
			return err
		}

		if (kind == reflect.Slice || kind == reflect.Array) && tp.Elem().Kind() != reflect.Uint8 {
			if len(e.nodes) == 0 {
				return e.endXML(p)
			}

			nod := e.nodes[0]

			for i, n := 0, val.Len(); i < n; i++ {
				vv := val.Index(i)

				err = nod.encodeXML(p, globVal, vv, true)
				if err != nil {
					return err
				}
			}

			return e.endXML(p)
		}

		if kind == reflect.Ptr {
			val = val.Elem()
		}

		if kind == reflect.Struct {
			_, ok := val.Interface().(time.Time)
			if ok {
				return errors.New("bind can't be time.Time")
			}

			for _, nn := range e.nodes {
				err = nn.encodeXML(p, globVal, val, false)
				if err != nil {
					return err
				}
			}
		}

		return e.endXML(p)
	case isItem && inVal.IsValid() && !inVal.IsZero():
		err := e.startXML(p, inVal)
		if err != nil {
			return err
		}

		if len(e.nodes) > 0 {
			for _, nn := range e.nodes {
				err = nn.encodeXML(p, globVal, inVal, false)
				if err != nil {
					return err
				}
			}
		} else {
			fieldVal := inVal.FieldByName(e.from)
			if fieldVal.Kind() == reflect.Ptr {
				fieldVal = fieldVal.Elem()
			}

			if fieldVal.Kind() == reflect.Map {
				return errors.New("field val must be only simple type")
			}

			if fieldVal.Kind() == reflect.Struct {
				iVal := fieldVal.Interface()

				t, ok := iVal.(time.Time)
				if !ok {
					return errors.New("static val must be only simple type")
				}

				_, err = p.WriteString(t.String())
				if err != nil {
					return err
				}
			} else {
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
			}
		}

		return e.endXML(p)
	case e.from != "" && inVal.IsValid() && !inVal.IsZero():
		err := e.startXML(p, inVal)
		if err != nil {
			return err
		}

		fieldVal := inVal.FieldByName(e.from)
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		if fieldVal.Kind() == reflect.Map {
			return errors.New("field val must be only simple type")
		}

		if fieldVal.Kind() == reflect.Struct {
			iVal := fieldVal.Interface()

			t, ok := iVal.(time.Time)
			if !ok {
				return errors.New("static val must be only simple type")
			}

			_, err = p.WriteString(t.String())
			if err != nil {
				return err
			}
		} else {
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
		}

		return e.endXML(p)
	default:
		zeroVal := reflect.Zero(reflect.TypeOf(0))

		err := e.startXML(p, zeroVal)
		if err != nil {
			return err
		}

		for _, nn := range e.nodes {
			err = nn.encodeXML(p, globVal, zeroVal, false)
			if err != nil {
				return err
			}

		}

		return e.endXML(p)
	}
}

func (e *node) startXML(p *printer, val reflect.Value) error {
	if strings.TrimSpace(e.to) == "" {
		return errors.New("node `to` field is empty")
	}

	_, err := p.WriteRune('<')
	if err != nil {
		return err
	}

	_, err = p.WriteString(e.to)
	if err != nil {
		return err
	}

	for _, a := range e.args {
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

func (e *node) endXML(p *printer) error {
	if strings.TrimSpace(e.to) == "" {
		return errors.New("node `to` field is empty")
	}

	_, err := p.WriteString("</")
	if err != nil {
		return err
	}

	_, err = p.WriteString(e.to)
	if err != nil {
		return err
	}

	_, err = p.WriteRune('>')
	if err != nil {
		return err
	}

	return nil
}

func (a *arg) encodeXML(p *printer, val reflect.Value) error {
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
			return errors.New("static val must be only simple type")
		}

		if kind == reflect.Struct {
			t, ok := statVal.Interface().(time.Time)
			if !ok {
				return errors.New("static val must be only simple type")
			}

			_, err = p.WriteString(t.String())
			if err != nil {
				return err
			}
		} else {
			s, b, err := marshalSimpleVal(statVal)
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
		}
	} else {
		kind := val.Kind()

		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
		}

		var fieldVal reflect.Value
		if kind == reflect.Struct {
			fieldVal = val.FieldByName(a.from)
			if !fieldVal.IsValid() {
				return errors.New("invalid field val for attr")
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
	}

	_, err = p.WriteString(`"`)
	if err != nil {
		return err
	}

	return nil
}
