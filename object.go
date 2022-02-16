package v2

import (
	"fmt"
	"strings"
)

//Object JSON object
type Object struct {
	Properties []*Property
	fields     map[string]int
	params     bool
	text       string
}

//New constucts Object
func New(fields ...Property) *Object {
	var jo Object
	if len(fields) > 0 {
		for _, f := range fields {
			jo.Add(f.Name, f.Value)
		}
	}
	return &jo
}

//ParseObject parses JSON object
func ParseObject(data []byte) (*Object, error) {
	return parseObject(data, false, false)
}

//ParseObjectSafe parses JSON object ignoring prefix non-ASCII characters
func ParseObjectSafe(data []byte) (*Object, error) {
	return parseObject(data, false, true)
}

//ParseObjectWithParameters parses parameterized JSON object
func ParseObjectWithParameters(data []byte) (*Object, error) {
	return parseObject(data, true, false)
}

func parseObject(data []byte, parameterized bool, safe bool) (*Object, error) {
	var (
		p   = newParser(data)
		err error
	)

	if safe {
		err = p.EnsureJSON()
	} else {
		err = p.SkipWS()
	}
	if err != nil {
		return nil, err
	}
	if p.Byte != '{' {
		return nil, fmt.Errorf("parsing JSON object; expect '{', found '%s'", string(p.Byte))
	}
	return p.ParseObject(parameterized)
}

//Add adds field to JSON object; if field exists overrides it; we expect valid name and value
func (jo *Object) Add(field string, value Value) {
	if value.Type() == ObjectType && jo == value {
		value, _ = copyValue(value)
	}
	jo.addProperty(&Property{
		Name:  field,
		Value: value,
	})
}

func (jo *Object) addProperty(jp *Property) {
	jo.text = ""

	if jo.fields == nil {
		jo.fields = map[string]int{jp.Name: 0}
		jo.Properties = []*Property{jp}
		return
	}

	i, ok := jo.fields[jp.Name]
	if ok {
		jo.Properties[i] = jp
		return
	}

	jo.fields[jp.Name] = len(jo.Properties)
	jo.Properties = append(jo.Properties, jp)
}

//Remove removes field from JSON object; ignore is field doesn't exist; we expect good name
func (jo *Object) Remove(field string) {
	jo.text = ""

	sz := len(jo.fields)
	if sz == 0 {
		return
	}

	if _, ok := jo.fields[field]; !ok {
		return
	}

	if sz == 1 {
		jo.fields = nil
		jo.Properties = nil
		return
	}

	var (
		fs     = make([]*Property, sz-1)
		fields = make(map[string]int)
		i      int
	)

	for _, jp := range jo.Properties {
		if jp.Name == field {
			continue
		}
		fs[i] = jp
		fields[jp.Name] = i
		i++
	}

	jo.Properties = fs
	jo.fields = fields
}

//GetProperty gets Property by field name
func (jo *Object) GetProperty(field string) (*Property, bool) {
	if len(jo.fields) == 0 {
		return nil, false
	}

	i, ok := jo.fields[field]
	if !ok {
		return nil, false
	}

	return jo.Properties[i], true
}

//GetValue gets value by field name
func (jo *Object) GetValue(field string) (Value, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.Value, true
}

//GetString gets string value
func (jo *Object) GetString(field string) (string, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return "", false
	}
	return jp.GetString()
}

//GetStrings gets string values
func (jo *Object) GetStrings(field string) ([]string, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetStrings()
}

//GetInt gets int value
func (jo *Object) GetInt(field string) (int, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return 0, false
	}
	return jp.GetInt()
}

//GetInts gets int values
func (jo *Object) GetInts(field string) ([]int, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetInts()
}

//GetFloat gets float64 value
func (jo *Object) GetFloat(field string) (float64, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return 0, false
	}
	return jp.GetFloat()
}

//GetFloats gets float64 values
func (jo *Object) GetFloats(field string) ([]float64, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetFloats()
}

//GetBool gets bool value
func (jo *Object) GetBool(field string) (bool, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return false, false
	}
	return jp.GetBool()
}

//GetObject gets Object value
func (jo *Object) GetObject(field string) (*Object, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetObject()
}

//GetObjects gets Object values
func (jo *Object) GetObjects(field string) ([]*Object, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetObjects()
}

//GetArray gets Object value
func (jo *Object) GetArray(field string) (*Array, bool) {
	jp, ok := jo.GetProperty(field)
	if !ok {
		return nil, false
	}
	return jp.GetArray()
}

//Copy copies Object
func (jo *Object) Copy() *Object {
	if jo == nil {
		return nil
	}

	if len(jo.Properties) == 0 {
		return new(Object)
	}

	copy := Object{
		params: jo.params,
		fields: make(map[string]int),
	}

	for _, jp := range jo.Properties {
		v, ok := copyValue(jp.Value)
		if ok {
			copy.fields[jp.Name] = len(copy.Properties)
			copy.Properties = append(copy.Properties, &Property{
				Name:  jp.Name,
				Value: v,
			})
		}
	}

	return &copy
}

func (jo *Object) IsEmpty() bool {
	return jo == nil || len(jo.Properties) == 0
}

//Equals compares two JSON objects
func (jo *Object) Equals(other *Object) (bool, error) {
	var (
		left  = jo
		right = other
	)

	if left == nil {
		left = New()
	}

	if right == nil {
		right = New()
	}

	for f := range left.fields {
		if _, ok := right.fields[f]; !ok {
			v, _ := left.GetValue(f)
			if !v.IsEmpty() {
				return false, fmt.Errorf("extra property: %s", f)
			}
		}
	}

	for f := range right.fields {
		if _, ok := left.fields[f]; !ok {
			v, _ := right.GetValue(f)
			if !v.IsEmpty() {
				return false, fmt.Errorf("missing property: %s", f)
			}
		}
	}

	for _, l := range left.Properties {
		for _, r := range right.Properties {
			if r.Name == l.Name {
				err := compareValues(l.Value, r.Value)
				if err != nil {
					return false, fmt.Errorf("mismatch property %s: %s", l.Name, err.Error())
				}
				break
			}
		}
	}

	return true, nil
}

//SetParameters replaces parameter placeholders with values
func (jo *Object) SetParameters(params *Object) *Object {
	var set Object

	for _, jp := range jo.Properties {
		var (
			name  = jp.Name
			value = jp.Value
		)
		if len(jp.namep) > 0 {
			name, _ = setStringParameters(fmt.Sprintf("\"%s\"", jp.Name), jp.namep, params)
			if len(name) < 3 {
				continue
			}
			name = string(name[1 : len(name)-1])
		}

		value = setValueParameters(value, jp.valuep, params)
		if value == nil {
			continue
		}
		switch value.Type() {
		case StringType:
			if value.Value() == "" {
				continue
			}
		case ObjectType:
			o, ok := value.(*Object)
			if !ok || len(o.Properties) == 0 {
				continue
			}
		case ArrayType:
			a, ok := value.(*Array)
			if !ok || len(a.Values) == 0 {
				continue
			}
		}
		set.Add(name, value)
	}

	return &set
}

//GetParameters retrieves paramaters from Object
func (jo *Object) GetParameters() []Parameter {
	var (
		params []Parameter
	)

	for _, jp := range jo.Properties {
		if len(jp.namep) > 0 {
			params = append(params, jp.namep...)
		}

		switch jp.Value.Type() {
		case StringType:
			if len(jp.valuep) > 0 {
				params = append(params, jp.valuep...)
			}
		case ObjectType:
			o, ok := jp.Value.(*Object)
			if ok {
				params = append(params, o.GetParameters()...)
			}
		case ArrayType:
			a, ok := jp.Value.(*Array)
			if ok {
				params = append(params, a.GetParameters()...)
			}
		}
	}

	return params
}

//Value Value interface
func (jo *Object) Value() interface{} {
	return jo.Properties
}

//Type Value interface
func (jo *Object) Type() ValueType {
	return ObjectType
}

//String Value interface
func (jo *Object) String() string {
	if jo == nil {
		return "{}"
	}

	if jo.text == "" {
		sz := len(jo.Properties)
		if sz == 0 {
			jo.text = "{}"
		} else {
			values := make([]string, sz)
			for i, jp := range jo.Properties {
				values[i] = fmt.Sprintf("\"%s\":%s", jp.Name, jp.Value.String())
			}
			jo.text = fmt.Sprintf("{%s}", strings.Join(values, ","))
		}
	}
	return jo.text
}

//used when a first byte is '{'
func (p *byteParser) ParseObject(parameterized bool) (*Object, error) {
	var jo Object

	for {
		err := p.SkipWS()
		if err != nil {
			return nil, err
		}
		if p.Byte != '"' {
			if p.Byte == '}' {
				break
			}
			return nil, fmt.Errorf("parsing object at %d: expected [ \" ], found %s", p.Index, string(p.Byte))
		}
		jp, err := p.ParseProperty(parameterized)
		if err != nil {
			return nil, err
		}

		if len(jp.valuep) > 0 || len(jp.namep) > 0 {
			jo.params = true
		}

		if p.Byte == ',' {
			jo.addProperty(jp)
			continue
		}

		if p.Byte == '}' {
			jo.addProperty(jp)
			break
		}
	}

	err := p.SkipWS()
	if err == errEOF {
		err = nil
	}

	return &jo, err
}
