package main

import (
	"bytes"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/tools/imports"
)

type Table struct {
	Name    string
	Comment string
	Columns []*Column
}

type Column struct {
	Name     string
	Comment  string
	DataType string
	Nullable bool
}

func (t *Table) ToStruct() ([]byte, error) {
	var b bytes.Buffer
	funcs := template.FuncMap{
		"TitleCase": TitleCase,
		"CamelCase": CamelCase,
		"DataType":  DataType,
	}

	templateBox, err := rice.FindBox("template")
	data, err := templateBox.Bytes("struct.text")
	if err != nil {
		return nil, err
	}
	tp, err := template.New("struct").Funcs(funcs).Parse(string(data))
	if err != nil {
		return nil, err
	}

	err = tp.Execute(&b, t)
	if err != nil {
		return nil, err
	}

	data, err = imports.Process("", b.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
