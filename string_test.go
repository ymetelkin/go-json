package v2

import (
	"fmt"
	"testing"
)

func TestStringParse(t *testing.T) {
	s := `"abc"`
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, _, err := p.ParseString(false)
	if err != errEOF {
		t.Error("Failed to parse [ abc ]")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `{ s: "\"YM\" is Yuri Metelkin's \stuff\ \/ \n initials" }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //s
	p.Read()         //:
	err = p.SkipWS() //"
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, params, err := p.ParseString(false)
	if err != nil {
		t.Error(err.Error())
	}
	if len(params) > 0 {
		t.Error("Failed to parse parameters")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `"\\\\"`
	p = newParser([]byte(s))
	p.SkipWS()
	_, _, err = p.ParseString(false)
	if err != errEOF {
		t.Error("Failed to parse [ \\ ]")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `{ s: "is \"concerned\\\" (second reference)\\" }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //s
	p.Read()         //:
	err = p.SkipWS() //"
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, params, err = p.ParseString(false)
	if err != nil {
		t.Error(err.Error())
	}
	if len(params) > 0 {
		t.Error("Failed to parse parameters")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `{ s: "\u0059\u004D" }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //s
	p.Read()         //:
	err = p.SkipWS() //"
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, params, err = p.ParseString(false)
	if err != nil {
		t.Error(err.Error())
	}
	if len(params) > 0 {
		t.Error("Failed to parse parameters")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `{ s: "value = ${value}" }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //s
	p.Read()         //:
	err = p.SkipWS() //"
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, params, err = p.ParseString(true)
	if err != nil {
		t.Error(err.Error())
	}
	if len(params) == 0 {
		t.Error("Failed to parse parameters")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())

	s = `{ s: "\"MSFT\" is $9.99" }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //s
	p.Read()         //:
	err = p.SkipWS() //"
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '"' {
		t.Error("Failed to parse [ \" ]")
	}
	v, params, err = p.ParseString(false)
	if err != nil {
		t.Error(err.Error())
	}
	if len(params) > 0 {
		t.Error("Failed to parse parameters")
	}
	fmt.Printf("%s\t%s\n", v.Value(), v.String())
}
