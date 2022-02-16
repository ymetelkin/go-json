package v2

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestObjectParse(t *testing.T) {
	s := `{ "text": "abc", "number": 3.14, "flag": true, "array": [ 1, 2, 3 ], "object": { "a": "b" }}`
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '{' {
		t.Error("Failed to parse {")
	}
	v, err := p.ParseObject(false)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(v.String())

	s = `{"x": "Arts &amp; Entertainment; a &lt; b or c &gt; d; YM & &Co"}`
	p = newParser([]byte(s))
	p.SkipWS()
	v, err = p.ParseObject(false)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v.String())
}

func TestEnsureJSON(t *testing.T) {
	s := `  ï » ¿{ "text": "abc"} ï » ¿ `
	jo, err := ParseObjectSafe([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(jo.String())

	s = ` ï » ¿<text>abc</text>`
	_, err = ParseObjectSafe([]byte(s))
	if err == nil {
		t.Error("Must not parse it")
	}
	fmt.Println(err.Error())

	s = `{ "text": "abc"]`
	_, err = ParseObjectSafe([]byte(s))
	if err == nil {
		t.Error("Must not parse it")
	}
	fmt.Println(err.Error())

	s = `  ï » ¿[{ "text": "abc"}] ï » ¿ `
	ja, err := ParseArraySafe([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(ja.String())
}

func TestObjectCopy(t *testing.T) {
	jo := New(
		Field("name", String("YM")),
		Field("null", Null()),
		Field("empty", String("")),
	)
	fmt.Println(jo.String())
	copy := jo.Copy()
	fmt.Println(copy.String())
}

func TestObjectPointers(t *testing.T) {
	jo := New(Field("name", String("YM")))
	fmt.Println(jo.String())
	jo.Add("person", jo)
	fmt.Println(jo.String())
}

func TestGraph(t *testing.T) {
	data, _ := ioutil.ReadFile("test_data/graph.json")
	jo, _ := ParseObject(data)
	ja, _ := jo.GetObjects("vertices")
	vertices := make(map[int]graphPerson)
	for i, v := range ja {
		name, _ := v.GetString("term")
		vertices[i] = graphPerson{
			Name: name,
		}
	}

	ja, _ = jo.GetObjects("connections")
	for _, o := range ja {
		source, _ := o.GetInt("source")
		target, _ := o.GetInt("target")
		weight, _ := o.GetFloat("weight")
		count, _ := o.GetInt("doc_count")
		v := vertices[source]
		c := vertices[target]
		v.Connections = append(v.Connections, graphConnection{
			Name:   c.Name,
			Weight: weight,
			Count:  count,
		})
		vertices[source] = v
	}

	for _, p := range vertices {
		if len(p.Connections) == 0 {
			continue
		}

		fmt.Println(p.Name)
		for _, c := range p.Connections {
			fmt.Printf("\t%s\n", c.Name)
		}
		fmt.Println()
	}
}

type graphPerson struct {
	Name        string
	Connections []graphConnection
}

type graphConnection struct {
	Name   string
	Weight float64
	Count  int
}
