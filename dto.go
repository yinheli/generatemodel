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
	Name          string
	TitleCaseName string
	Comment       string
	Columns       []*Column
}

type Column struct {
	Name          string
	Comment       string
	DataType      string
	Nullable      bool
	TitleCaseName string
	CamelCaseName string
	GoType        string
	Tag           string
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
		l.Println("origin data: ", string(b.Bytes()))
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

	structTemplate, err := template.New("struct").Parse(string(data))
	if err != nil {
		return nil, err
	}

	return structTemplate, nil
}

func Tag(column Column) string {
	jsonTag := fmt.Sprintf(`json:"%s"`, JsonTag(column.CamelCaseName, column.GoType))
	validateTag := ""
	if !column.Nullable {
		switch column.TitleCaseName {
		case "Id", "CreatedAt", "UpdatedAt":
		default:
			validateTag = ` validate:"required"`
		}
	}

	return fmt.Sprintf("`%s%s`", jsonTag, validateTag)
}

func JsonTag(colunm string, goType string) string {
	switch goType {
	case "uint32", "int64", "uint64",
		"*uint32", "*int64", "*uint64":
		colunm = fmt.Sprintf("%s,string", colunm)
	}
	return colunm
}
