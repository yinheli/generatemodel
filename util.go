package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
	numberReplacement = []byte(`$1 $2 $3`)
)

func TitleCase(str string) string {
	return toCamelCase(str, true)
}

func CamelCase(str string) string {
	return toCamelCase(str, false)
}

func DataType(dataType string, nullable bool) string {
	dataType = strings.ToLower(strings.TrimSpace(dataType))

	goType := "string"

	switch strings.TrimSpace(dataType) {
	case "int", "tinyint":
		goType = "int64"
	case "date", "datetime", "timestamp":
		goType = "time.Time"
	case "float", "decimal", "double":
		goType = "float64"
	}

	if nullable {
		return fmt.Sprintf("*%s", goType)
	}

	return goType
}

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// Converts a string to CamelCase
func toCamelCase(s string, first bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := first
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}
