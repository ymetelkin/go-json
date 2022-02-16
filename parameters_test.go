package v2

import (
	"fmt"
	"testing"
)

func TestParameters(t *testing.T) {
	s := `"${name}"`
	p := newParser([]byte(s))
	p.SkipWS()
	_, pms, _ := p.ParseString(true)
	if len(pms) != 1 {
		t.Error("must be one parameter")
	}
	if pms[0].Name != "name" {
		fmt.Println(pms[0].Name)
		t.Error("Failed to parse parameter name")
	}
	if pms[0].Value.String() != `""` {
		fmt.Println(pms[0].Value.String())
		t.Error("Failed to parse parameter value")
	}
	params, _ := ParseObject([]byte(`{"name":"YM"}`))
	s2, v := setStringParameters(s, pms, params)
	if s2 != `"YM"` || v.Type() != StringType || v.String() != `"YM"` {
		t.Error("Failed to set string parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)

	s = `"${name?YM}"`
	p = newParser([]byte(s))
	p.SkipWS()
	_, pms, _ = p.ParseString(true)
	if len(pms) != 1 {
		t.Error("must be one parameter")
	}
	if pms[0].Name != "name" {
		fmt.Println(pms[0].Name)
		t.Error("Failed to parse parameter name")
	}
	if pms[0].Value.String() != `"YM"` {
		fmt.Println(pms[0].Value.String())
		t.Error("Failed to parse parameter value")
	}
	s2, v = setStringParameters(s, pms, nil)
	if s2 != `"YM"` || v.Type() != StringType || v.String() != `"YM"` {
		t.Error("Failed to set string parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)
	params, _ = ParseObject([]byte(`{"name":"SV"}`))
	s2, v = setStringParameters(s, pms, params)
	if s2 != `"SV"` || v.Type() != StringType || v.String() != `"SV"` {
		t.Error("Failed to set string parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)

	s = `"n${a}${m?m}e${num}"`
	p = newParser([]byte(s))
	p.SkipWS()
	_, pms, _ = p.ParseString(true)
	if len(pms) != 3 {
		t.Error("must be 3 parameters")
	}
	if pms[0].Name != "a" {
		fmt.Println(pms[0].Name)
		t.Error("Failed to parse parameter name")
	}
	if pms[1].Value.String() != `"m"` {
		fmt.Println(pms[0].Value.String())
		t.Error("Failed to parse parameter value")
	}
	params, _ = ParseObject([]byte(`{"a":"a","num":1}`))
	s2, v = setStringParameters(s, pms, params)
	if s2 != `"name1"` || v.Type() != StringType || v.String() != `"name1"` {
		fmt.Println(v.String())
		t.Error("Failed to set string parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)

	s = `"${f:float?3.14}"`
	p = newParser([]byte(s))
	p.SkipWS()
	_, pms, _ = p.ParseString(true)
	if len(pms) != 1 {
		t.Error("must be one parameter")
	}
	if pms[0].Name != "f" {
		fmt.Println(pms[0].Name)
		t.Error("Failed to parse parameter name")
	}
	if pms[0].Value.String() != "3.14" {
		fmt.Println(pms[0].Value.String())
		t.Error("Failed to parse parameter value")
	}
	s2, v = setStringParameters(s, pms, nil)
	if s2 != `3.14` || v.Type() != FloatType || v.Value() != 3.14 {
		fmt.Println(v.String())
		t.Error("Failed to set float parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)
	params, _ = ParseObject([]byte(`{"f":0.666}`))
	s2, v = setStringParameters(s, pms, params)
	if s2 != "0.666" || v.Type() != FloatType || v.Value() != 0.666 {
		fmt.Println(v.Value())
		t.Error("Failed to set float parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)

	s = `"${o:object?{\"name\":\"YM\"}}"`
	p = newParser([]byte(s))
	p.SkipWS()
	_, pms, _ = p.ParseString(true)
	if len(pms) != 1 {
		t.Error("must be one parameter")
	}
	if pms[0].Name != "o" {
		fmt.Println(pms[0].Name)
		t.Error("Failed to parse parameter name")
	}
	if pms[0].Value.String() != `{"name":"YM"}` {
		fmt.Println(pms[0].Value.String())
		t.Error("Failed to parse parameter value")
	}
	s2, v = setStringParameters(s, pms, nil)
	if s2 != `{"name":"YM"}` || v.Type() != ObjectType || v.String() != `{"name":"YM"}` {
		fmt.Println(v.String())
		t.Error("Failed to set object parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)
	params, _ = ParseObject([]byte(`{"o":{"name":"SV"}}`))
	s2, v = setStringParameters(s, pms, params)
	if s2 != `{"name":"SV"}` || v.Type() != ObjectType || v.String() != `{"name":"SV"}` {
		fmt.Println(v.String())
		t.Error("Failed to set object parameters")
	}
	fmt.Printf("%s -> %s\n", s, s2)
}

func TestObjectParameters(t *testing.T) {
	s := `{ "name": "${name?YM}", "age": "${age:int}", "salary": "${salary:float}", "url":"${domain?www.ap.org}/${action?story}/${id}", "a${p1?b}${p2?d}":123}`
	jo, err := ParseObjectWithParameters([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	params, _ := ParseObject([]byte(`{"name":"SV", "age":27, "salary": 80000, "id":"xyz", "p2":"c"}`))
	jo = jo.SetParameters(params)
	fmt.Println(jo.String())

	s = `{ "index": "${index?appl}", "query": "${query:object}", "size": "${size:int?10}", "sort": "${sort:array}"}`
	jo, err = ParseObjectWithParameters([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	params, _ = ParseObject([]byte(`{"index":"appl-thirty", "query":{"query_string":{"query":"test"}}, "size":40, "sort": [{"arrivaldatetime":"desc"}]}`))
	jo = jo.SetParameters(params)
	fmt.Println(jo.String())

	s = `{ "index": "${index?appl}", "query": {"query_string":{"query":"${query}","fields":"${fields:array}"}}, "size": "${size:int?10}", "sort": [{"${sort_field}":"${sort_direction}"}]}`
	jo, err = ParseObjectWithParameters([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	params, _ = ParseObject([]byte(`{"query":"test", "fields":["type","headline"], "sort_field": "arrivaldatetime", "sort_direction":"desc"}`))
	jo = jo.SetParameters(params)
	fmt.Println(jo.String())
}

func TestV1Parsing(t *testing.T) {
	input := `{"id":"${id}","name":"${name}"}`
	expected := `{"id":1,"name":"YM"}`
	jo, err := ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms := jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ := ParseObject([]byte(`{"id":1,"name":"YM"}`))
	jo = jo.SetParameters(params)
	test := jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"id":"${id}", "name":"${name}"}`
	expected = `{"name":"YM"}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"name":"YM"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"id":"${id}", "name":"${name}","products":["${p1}","${p2}"]}`
	expected = `{"id":1,"name":"YM","products":[1,2]}`
	jo, err = ParseObjectWithParameters([]byte(input))
	fmt.Println(jo.String())
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"id":1,"name":"YM","index":"appl","size":20,"p1":1,"p2":2}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"id":"${id}","name":"${name}", "child":{"name":"${name}", "age":"${age}","${extra_field}":"${extra_value}"}}`
	expected = `{"id":1,"name":"YM","child":{"name":"YM","age":13,"nick":"Gusyonok"}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"id":1,"name":"YM","age":13,"extra_field":"nick","extra_value":"Gusyonok"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"id":"${id}","name":"${name}", "child":{"age":"${age}"}}`
	expected = `{"id":1,"name":"YM"}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"id":1,"name":"YM"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"${id_prefix}_id":"${id}","name":"${name}", "child":{"name":"${name} Jr.", "age":"${age}"}}
`
	expected = `{"user_id":1,"name":"YM","child":{"name":"YM Jr."}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"id_prefix":"user","id":1,"name":"YM"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"${id_prefix}_id":"${id}","name":"${name}", "child":{"name":"${name} Jr.", "age":"${age}"}}`
	expected = `{"_id":1,"name":"YM","child":{"name":"YM Jr."}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"id_prefix1":"user","id":1,"name":"YM"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"template":{"query":{"query_string":{"query":"${query}","fields":"${fields}"}}}}`
	expected = `{"template":{"query":{"query_string":{"query":"test","fields":["head","body"]}}}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"query":"test","fields":["head","body"]}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"template":{"_source":["headline","${field1}","${field2}","${field3}"]}}`
	expected = `{"template":{"_source":["headline","type","date"]}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"field1":"type","field2":"date"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"first_name":"${first_name?Yuri}","last_name":"${last_name?Metelkine}"}`
	expected = `{"first_name":"Yuri","last_name":"Metelkin"}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"last_name":"Metelkin"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"name":"${first_name?Yuri} ${last_name?Metelkine}"}`
	expected = `{"name":"Yuri Metelkin"}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"last_name":"Metelkin"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"template":{"_source":["headline","${field1}","${field2}","${field3?test}"]}}`
	expected = `{"template":{"_source":["headline","type","date","test"]}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"field1":"type","field2":"date"}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"query":{"bool":{"must":{"query_string":{"query":"${query}","fields":["head","body"]}},"filter":[{"terms":{"type":"${media_types}"}},{"terms":{"filing.products":"${include_products}"}}],"must_not":{"terms":{"filing.products":"${exclude_products}"}}}}}`
	expected = `{"query":{"bool":{"must":{"query_string":{"query":"ap","fields":["head","body"]}},"filter":[{"terms":{"type":["audio"]}},{"terms":{"filing.products":[1,2]}}],"must_not":{"terms":{"filing.products":[3]}}}}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	pms = jo.GetParameters()
	for _, p := range pms {
		fmt.Printf("%s\t%s\n", p.Name, p.Value.String())
	}
	params, _ = ParseObject([]byte(`{"query":"ap","media_types":["audio"],"include_products":[1,2], "exclude_products":[3]}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}

	input = `{"query":{"bool":{"must":{"query_string":{"query":"${query}","fields":["head","body"]}},"filter":[{"terms":{"type":"${media_types}"}},{"terms":{"filing.products":"${include_products}"}}],"must_not":{"terms":{"filing.products":"${exclude_products}"}}}}}`
	expected = `{"query":{"bool":{"must":{"query_string":{"query":"ap","fields":["head","body"]}},"filter":[{"terms":{"type":["audio"]}}]}}}`
	jo, err = ParseObjectWithParameters([]byte(input))
	if err != nil {
		t.Error(err.Error())
	}
	params, _ = ParseObject([]byte(`{"query":"ap","media_types":["audio"]}`))
	jo = jo.SetParameters(params)
	test = jo.String()
	fmt.Println(test)
	if test != expected {
		t.Error("Doesn't match!")
		fmt.Println(test)
	}
}
