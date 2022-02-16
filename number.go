package v2

import (
	"fmt"
	"strconv"
	"strings"
)

//returns float64 or int
func (p *byteParser) ParseNumber() (Value, error) {
	var (
		sb    strings.Builder
		float bool
		pos   = p.Byte == '+'
		sign  = pos || p.Byte == '-'
		idx   = p.Index
		err   error
	)

	if !pos {
		sb.WriteByte(p.Byte)
	}

	if sign && isWS(p.Byte) {
		err = p.SkipWS()
		if err != nil {
			return nil, fmt.Errorf("parsing number at %d: %s", idx, err.Error())
		}
	}

	for {
		err = p.Read()
		if err != nil {
			if err == errEOF {
				break
			}
			return nil, fmt.Errorf("parsing number at %d: %s", idx, err.Error())
		}

		if p.Byte >= '0' && p.Byte <= '9' {
			sb.WriteByte(p.Byte)
		} else if p.Byte == '.' || p.Byte == 'e' || p.Byte == '-' || p.Byte == '+' {
			float = true
			sb.WriteByte(p.Byte)
		} else if p.Byte == 'E' {
			float = true
			sb.WriteByte('e')
		} else {
			break
		}
	}

	var (
		s  = sb.String()
		jv Value
	)

	if float {
		v, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return nil, fmt.Errorf("parsing number at %d: %s", idx, e.Error())
		}
		jv = &floatValue{
			value: v,
			text:  s,
		}
	} else {
		v, e := strconv.Atoi(s)
		if e != nil {
			if strings.Contains(e.Error(), "value out of range") && v > 0 {
				var u uint64
				u, e = strconv.ParseUint(s, 10, 0)
				if e == nil {
					jv = &uintValue{
						value: u,
						text:  s,
					}
				}
			}
			if e != nil {
				return nil, fmt.Errorf("parsing number at %d: %s", idx, e.Error())
			}
		} else {
			jv = &intValue{
				value: v,
				text:  s,
			}
		}
	}

	if isWS(p.Byte) && err == nil {
		err = p.SkipWS()
	}

	return jv, err
}
