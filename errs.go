package gotemplconstr

import "errors"

var (
	ErrToFiledEmpty            = errors.New("node `to` field is empty")
	ErrStaticValOnlySimpleType = errors.New("static value field must be only simple type")
	ErrValOnlySimpleType       = errors.New("value must be only simple type")
	ErrBindCantTime            = errors.New("bind can't be time.Time type")
	ErrInvalidField            = errors.New("invalid binding `from` field")
)
