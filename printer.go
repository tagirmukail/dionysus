package gotemplconstr

import (
	"bytes"
	"io"
	"reflect"
)

const flushSize = 1024 // bytes

type printer struct {
	io.Writer

	tmp *bytes.Buffer
}

func (p printer) WriteByte(c byte) (err error) {
	defer func() {
		if err == nil {
			err = p.flush(false)
		}
	}()

	return p.tmp.WriteByte(c)
}

func (p printer) WriteString(s string) (n int, err error) {
	defer func() {
		if err == nil {
			err = p.flush(false)
		}
	}()

	return p.tmp.WriteString(s)
}

func (p printer) WriteRune(r rune) (n int, err error) {
	defer func() {
		if err == nil {
			err = p.flush(false)
		}
	}()

	return p.tmp.WriteRune(r)
}

func (p printer) writeNewLine() error {
	_, err := p.WriteRune('\n')

	return err
}

func (p printer) writeQuotesString(s string) (err error) {
	_, err = p.WriteRune('"')
	if err != nil {
		return err
	}

	_, err = p.WriteString(s)
	if err != nil {
		return err
	}

	_, err = p.WriteRune('"')
	if err != nil {
		return err
	}

	return nil
}

func (p printer) writeQuotesReflectString(kind reflect.Kind, s string) (err error) {
	if kind == reflect.String {
		_, err = p.WriteRune('"')
		if err != nil {
			return err
		}
	}

	_, err = p.WriteString(s)
	if err != nil {
		return err
	}

	if kind == reflect.String {
		_, err = p.WriteRune('"')
		if err != nil {
			return err
		}
	}

	return nil
}

func (p printer) flush(finish bool) error {
	if p.tmp.Len() < flushSize && !finish {
		return nil
	}

	_, err := p.Write(p.tmp.Bytes())
	if err != nil {
		return err
	}

	p.tmp.Reset()

	return nil
}
