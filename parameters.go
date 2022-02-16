package v2

import (
	"fmt"
	"strconv"
	"strings"
)

//Parameter is JSON parameter
//Expected format ${name:type?value} "type" is optional, default is "string"
type Parameter struct {
	Name  string
	Value Value
	Start int
	End   int
}

//used when a first byte is '$'; returns Parameter
//we expect correct format, so not much validation will be performed
//start is global index for the string start
func (p *byteParser) ParseParameter(start int) (Parameter, error) {
	var (
		param  Parameter
		val    string
		pcolon = 0
		pvalue = 0
		vt     = "string"
		curly  = 1
	)

	err := p.Read()
	if err != nil {
		return param, err
	}

	if p.Byte != '{' {
		//need to go back one byte
		p.Index--
		p.Byte = p.Bytes[p.Index]
		return param, nil
	}

	param.Start = p.Index - 1 - start

	for {
		err = p.Read()
		if err != nil {
			return param, err
		}

		if p.Byte == ':' && pcolon == 0 { //type
			pcolon = p.Index
		} else if p.Byte == '?' && pvalue == 0 {
			pvalue = p.Index
		} else if p.Byte == '}' {
			curly--
			if curly > 0 {
				continue
			}

			param.End = p.Index + 1 - start

			end := p.Index
			if pcolon > 0 {
				end = pcolon
			} else if pvalue > 0 {
				end = pvalue
			}

			param.Name = string(p.Bytes[param.Start+start+2 : end]) //start+2 "${

			if pvalue > 0 {
				val = string(p.Bytes[pvalue+1 : p.Index])
			}

			if pcolon > 0 {
				end := p.Index
				if pvalue > 0 {
					end = pvalue
				}
				vt = string(p.Bytes[pcolon+1 : end])
			}

			switch vt {
			case "string":
				param.Value = &stringValue{
					value: val,
				}
			case "int":
				v := intValue{}
				if val != "" {
					v.value, err = strconv.Atoi(val)
				}
				param.Value = &v
			case "float":
				v := floatValue{}
				if val != "" {
					v.value, err = strconv.ParseFloat(val, 64)
				}
				param.Value = &v
			case "bool":
				v := boolValue{}
				if val == "true" {
					v.value = true
				}
				param.Value = &v
			case "object":
				if val == "" {
					param.Value = &Object{}
				} else {
					val = strings.ReplaceAll(val, "\\\"", "\"")
					jo, e := ParseObject([]byte(val))
					param.Value = jo
					err = e
				}
			case "array":
				if val == "" {
					param.Value = &Array{}
				} else {
					val = strings.ReplaceAll(val, "\\\"", "\"")
					ja, e := ParseArray([]byte(val))
					param.Value = ja
					err = e
				}
			}

			return param, err
		} else if p.Byte == '{' {
			curly++
		}
	}
}

func setStringParameters(s string, pms []Parameter, params *Object) (string, Value) {
	sz := len(pms)
	if sz == 0 {
		val, err := strconv.Unquote(s)
		if err != nil {
			return s, nil
		}
		return s, &stringValue{
			value: val,
		}
	}

	if sz == 1 && pms[0].Start == 1 && pms[0].End == len(s)-1 {
		val := pms[0].Value
		if params != nil {
			v, ok := params.GetValue(pms[0].Name)
			if ok {
				val = v
			}
		}
		return val.String(), val
	}

	var (
		sb  strings.Builder
		end = 0
		bs  = []byte(s)
	)

	for _, pm := range pms {
		val := pm.Value
		if params != nil {
			v, ok := params.GetValue(pm.Name)
			if ok {
				val = v
			}
		}

		if pm.Start > end {
			sb.WriteString(string(bs[end:pm.Start]))
		}
		sb.WriteString(fmt.Sprintf("%v", val.Value()))
		end = pm.End
	}

	if end < len(s) {
		sb.WriteString(string(bs[end:]))
	}

	text := sb.String()
	value, _ := strconv.Unquote(text)
	return text, &stringValue{
		value: value,
	}
}

func setValueParameters(v Value, pms []Parameter, params *Object) Value {
	switch v.Type() {
	case StringType:
		_, value := setStringParameters(v.String(), pms, params)
		return value
	case ObjectType:
		jo, ok := v.(*Object)
		if ok {
			return jo.SetParameters(params)
		}
	case ArrayType:
		ja, ok := v.(*Array)
		if ok {
			return ja.SetParameters(params)
		}
	}
	return v
}
