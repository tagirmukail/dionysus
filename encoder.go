package dionysus

import (
	"fmt"
	"reflect"
)

func (t *Template) Encode(v interface{}) error {
	switch t.outputType {
	//case JSON:
	//	return t.encodeJSON(v)
	//case CSV, TSV:
	//	return t.encodeCSV(v)
	//case YML:
	//	return t.encodeYML(v)
	case XML:
		return t.encodeXML(reflect.ValueOf(v))
	default:
		return fmt.Errorf("unknown file type: %v", t.outputType)
	}
}
