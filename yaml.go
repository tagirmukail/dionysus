package gotemplconstr

import (
	"reflect"
	"time"
)

const defaultPadding = "  "

func (t *Template) encodeYAML(v reflect.Value) error {

	err := t.node.encodeYAML(&t.printer, 0, v, reflect.Zero(reflect.TypeOf(0)),
		false, true, false)
	if err != nil {
		return err
	}

	return t.printer.flush(true)
}

func (nn *node) encodeYAML(p *printer, padding uint, globVal, inVal reflect.Value,
	isItem, incrPadding, isItemField bool) (err error) {
	if nn == nil {
		return nil
	}

	if !isItemField {
		err = nn.yamlStart(p, padding, isItem, nn.staticVal != nil)
		if err != nil {
			return err
		}
	}

	if incrPadding {
		padding++
	}

	switch {
	case nn.staticVal != nil:
		val := reflect.ValueOf(nn.staticVal)

		kind := val.Kind()
		if kind == reflect.Ptr {
			val = val.Elem()
			kind = val.Kind()
		}

		if kind == reflect.Map {
			return ErrStaticValOnlySimpleType
		}

		if kind == reflect.Struct {
			iVal := val.Interface()

			t, ok := iVal.(time.Time)
			if !ok {
				return ErrStaticValOnlySimpleType
			}

			err = p.writeQuotesString(t.String())
			if err != nil {
				return err
			}

			err = p.writeNewLine()
			if err != nil {
				return err
			}

			break
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
			err = p.writeQuotesReflectString(kind, s)
			if err != nil {
				return err
			}
		}

		err = p.writeNewLine()
		if err != nil {
			return err
		}
	case nn.bind != "":
		val, ok := getVal(nn.bind, globVal)
		if !ok {
			return nil
		}

		kind := val.Kind()
		tp := val.Type()

		if (kind == reflect.Slice || kind == reflect.Array) && tp.Elem().Kind() != reflect.Uint8 {
			if len(nn.nodes) == 0 {
				break
			}

			nod := nn.nodes[0]

			for i, n := 0, val.Len(); i < n; i++ {
				vv := val.Index(i)

				err = nod.encodeYAML(p, padding, globVal, vv, true, true, false)
				if err != nil {
					return err
				}
			}
		}

		if kind == reflect.Ptr {
			val = val.Elem()
		}

		if kind != reflect.Struct {
			break
		}

		_, ok = val.Interface().(time.Time)
		if ok {
			return ErrBindCantTime
		}

		for _, n := range nn.nodes {
			err = n.encodeYAML(p, padding, globVal, val, false, true, false)
			if err != nil {
				return err
			}
		}
	case isItem && inVal.IsValid() && !inVal.IsZero():
		if len(nn.nodes) > 0 {
			err = writePadding(p, padding)
			if err != nil {
				return err
			}

			_, err = p.WriteString("- ")
			if err != nil {
				return err
			}

			for i, n := range nn.nodes {
				if i != 0 {
					err = writePadding(p, padding+1)
					if err != nil {
						return err
					}
				}

				_, err = p.WriteString(n.to)
				if err != nil {
					return err
				}

				_, err = p.WriteString(": ")
				if err != nil {
					return err
				}

				err = n.encodeYAML(p, padding, globVal, inVal, false, false, true)
				if err != nil {
					return err
				}
			}

			break
		}

		fieldVal := inVal.FieldByName(nn.from)
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		kind := fieldVal.Kind()

		if kind == reflect.Map {
			return ErrValOnlySimpleType
		}

		if kind == reflect.Struct {
			iVal := fieldVal.Interface()

			t, ok := iVal.(time.Time)
			if !ok {
				return ErrStaticValOnlySimpleType
			}

			_, err = p.WriteRune('"')
			if err != nil {
				return err
			}

			_, err = p.WriteString(t.String())
			if err != nil {
				return err
			}

			_, err = p.WriteRune('"')
			if err != nil {
				return err
			}

			break
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
			err = p.writeQuotesReflectString(kind, s)
			if err != nil {
				return err
			}
		}
	case nn.from != "" && inVal.IsValid() && !inVal.IsZero():
		if !isItemField {
			err = writePadding(p, padding)
			if err != nil {
				return err
			}
		}

		fieldVal := inVal.FieldByName(nn.from)
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		kind := fieldVal.Kind()

		if kind == reflect.Map {
			return ErrValOnlySimpleType
		}

		if kind == reflect.Struct {
			iVal := fieldVal.Interface()

			t, ok := iVal.(time.Time)
			if !ok {
				return ErrStaticValOnlySimpleType
			}

			err = p.writeQuotesString(t.String())
			if err != nil {
				return err
			}

			err = p.writeNewLine()
			if err != nil {
				return err
			}

			break
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
			err = p.writeQuotesReflectString(kind, s)
			if err != nil {
				return err
			}
		}

		err = p.writeNewLine()
		if err != nil {
			return err
		}
	default:
		for _, n := range nn.nodes {
			err = n.encodeYAML(p, padding, globVal, inVal, false, true, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (nn *node) yamlStart(p *printer, padding uint, isItem, isStatic bool) (err error) {
	if !isItem && padding > 0 {
		err = writePadding(p, padding)
		if err != nil {
			return err
		}
	}

	if isItem {
		return nil
	}

	_, err = p.WriteString(nn.to)
	if err != nil {
		return err
	}

	if isStatic {
		_, err = p.WriteString(": ")
		if err != nil {
			return err
		}
	} else {
		_, err = p.WriteString(":\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func writePadding(p *printer, padding uint) (err error) {
	for i := 0; i < int(padding); i++ {
		_, err = p.WriteString(defaultPadding)
		if err != nil {
			return err
		}
	}

	return nil
}
