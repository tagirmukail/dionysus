package dionysus

import (
	"bytes"
	"io"
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
