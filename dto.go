package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/tools/imports"
)

var (
	structTemplate *template.Template
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
	GoType   string
}

func (t *Table) ToStruct() ([]byte, error) {
	var b bytes.Buffer

	tp, err := getStructTemplate()
	if err != nil {
		return nil, err
	}

	err = tp.Execute(&b, t)
	if err != nil {
		return nil, err
	}

	data, err := imports.Process("", b.Bytes(), nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getStructTemplate() (*template.Template, error) {
	if structTemplate != nil {
		return structTemplate, nil
	}

	templateBox, err := rice.FindBox("template")
	data, err := templateBox.Bytes("struct.text")
	if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{
		"TitleCase": TitleCase,
		"CamelCase": CamelCase,
		"DataType":  DataType,
		"JsonTag":   JsonTag,
	}

	structTemplate, err := template.New("struct").Funcs(funcs).Parse(string(data))
	if err != nil {
		return nil, err
	}

	return structTemplate, nil
}

func JsonTag(colunm string, goType string) string {
	res := CamelCase(colunm)
	switch goType {
	case "uint32", "int64", "uint64",
		"*uint32", "*int64", "*uint64":
		res = fmt.Sprintf("%s,string", res)
	}
	return res
}
