package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelCase(t *testing.T) {
	data := map[string]string{
		"a_": "a",
		"app_name": "appName",
		"_app_name": "AppName",
	}

	for k, v := range data {
		assert.Equalf(t, v, CamelCase(k), "%v should %v", k, v)
	}
}

func TestTitleCase(t *testing.T) {
	data := map[string]string{
		"a_": "A",
		"app_name": "AppName",
		"_app_name": "AppName",
	}

	for k, v := range data {
		assert.Equalf(t, v, TitleCase(k), "%v should %v", k, v)
	}
}
