package v2

import (
	"errors"
	"fmt"
)

type byteParser struct {
	Bytes []byte
	Byte  byte
	Index int
	Size  int
}

var errEOF = errors.New("EOF")

func newParser(data []byte) byteParser {
	return byteParser{
		Bytes: data,
		Size:  len(data),
		Index: -1,
	}
}

func (p *byteParser) Read() error {
	p.Index++
	if p.Index == p.Size {
		p.Index--
		return errEOF
	}
	p.Byte = p.Bytes[p.Index]
	return nil
}

func (p *byteParser) SkipWS() error {
	for {
		err := p.Read()
		if err != nil {
			return err
		}
		if !isWS(p.Byte) {
			break
		}
	}

	return nil
}

func (p *byteParser) EnsureJSON() error {
	var (
		i, start int
		upd, ok  bool
		end      = len(p.Bytes)
		last     = byte('}')
	)

	if end == 0 {
		return nil
	}

	for i < end {
		c := p.Bytes[i]

		switch c {
		case '{':
			last = '}'
		case '[':
			last = ']'
		default:
			if c > 32 && c < 127 {
				return fmt.Errorf("JSON cannot start with '%s'", string(c))
			}
			i++
			continue
		}

		start = i
		upd = i > 0
		break
	}

	i = end - 1
	for i > 0 {
		c := p.Bytes[i]

		if c == last {
			i = i + 1
			upd = i < end
			end = i
			ok = start+1 < end
			break
		}

		if c > 32 && c < 127 {
			return fmt.Errorf("JSON cannot end with '%s'", string(c))
		}

		i--
	}

	if !ok {
		return fmt.Errorf("JSON must end with '%s'", string(last))
	}

	if upd {
		data := p.Bytes[start:end]
		p.Bytes = data
		p.Size = len(data)
		p.Index = -1
	}

	return p.Read()
}

func isWS(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r' || c == '\f' || c == '\v' || c == '\b'
}
