package v2

import (
	"fmt"
	"strconv"
	"strings"
)

type stringValue struct {
	value string
}

//String creates StringType value
func String(v string) Value {
	return &stringValue{
		value: v,
	}
}

func (v *stringValue) Value() interface{} {
	return v.value
}

func (v *stringValue) Type() ValueType {
	return StringType
}

func (v *stringValue) String() string {
	return strconv.Quote(v.value)
}

func (v *stringValue) Copy() Value {
	return String(v.value)
}

func (v *stringValue) IsEmpty() bool {
	return v == nil || v.value == ""
}

//used when a first byte is '"'; beside Value and error returns []Parameter
func (p *byteParser) ParseString(parameterized bool) (Value, []Parameter, error) {
	var (
		idx     = p.Index
		params  []Parameter
		escapes []int
	)
	for {
		err := p.Read()
		if err != nil {
			return nil, nil, fmt.Errorf("parsing string at %d: %s", idx, err.Error())
		}

		if p.Byte == '"' {
			if p.Bytes[p.Index-1] == '\\' {
				//escaped quote is ignored but double '\' is a literal "\"
				var (
					i   = p.Index - 2
					dbl bool
				)
				for i > idx {
					if p.Bytes[i] != '\\' {
						break
					}
					dbl = !dbl
					i--
				}
				if !dbl {
					continue
				}
			}
			if len(escapes) == 0 {
				s := string(p.Bytes[idx+1 : p.Index])
				return &stringValue{
					value: s,
				}, params, p.SkipWS()
			}

			var (
				sb   strings.Builder
				last = idx + 1
			)
			for _, i := range escapes {
				if i > last {
					sb.Write(p.Bytes[last:i])
				}
				if i < last {
					continue
				}

				c := p.Bytes[i+1]
				switch c {
				case '"', '\'', '\\', '/':
					sb.WriteByte(c)
				case 'n':
					sb.WriteByte('\n')
				case 'r':
					sb.WriteByte('\r')
				case 't':
					sb.WriteByte('\t')
				case 'a':
					sb.WriteByte('\a')
				case 'b':
					sb.WriteByte('\b')
				case 'f':
					sb.WriteByte('\f')
				case 'v':
					sb.WriteByte('\v')
				case 'x', 'u', 'U', '0', '1', '2', '3', '4', '5', '6', '7':
					//this is too much, let's GO handle it
					s, err := strconv.Unquote(string(p.Bytes[idx : p.Index+1]))
					if err != nil {
						return nil, nil, fmt.Errorf("parsing string at %d: %s", idx, err.Error())
					}
					return &stringValue{
						value: s,
					}, params, p.SkipWS()
				default:
					sb.WriteByte('\\')
					sb.WriteByte(c)
				}
				last = i + 2
			}

			if last < p.Index+1 {
				sb.Write(p.Bytes[last:p.Index])
			}

			return &stringValue{
				value: sb.String(),
			}, params, p.SkipWS()
		} else if p.Byte == '$' && parameterized {
			param, err := p.ParseParameter(idx)
			if param.Name != "" {
				params = append(params, param)
			}
			if err != nil {
				return nil, nil, fmt.Errorf("parsing string at %d: %s", idx, err.Error())
			}
		} else if p.Byte == '\\' { //escapes; requires special handling
			escapes = append(escapes, p.Index)
		}
	}
}
