package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Setenv("GOFILE", "/tmp/gen/tmp")
	os.Setenv("gm_overwrite", "true")
	os.Setenv("database_uri", "xpush:xpush@(192.168.4.23:3306)/xpush?charset=utf8&parseTime=True&loc=Local")
}

func Test_tables(t *testing.T) {
	err := openDB()
	assert.Nil(t, err)
	tables, err := tables()
	assert.Nil(t, err)

	if len(tables) > 0 {
		l.Printf("Name: %s Comment: %s", tables[0].Name, tables[1].Comment)
	}
}

func Test_columns(t *testing.T) {
	err := openDB()
	assert.Nil(t, err)
	tables, err := tables()
	assert.Nil(t, err)

	assert.Equal(t, len(tables) > 0, true)

	columns, err := columns(tables[0].Name)
	assert.Nil(t, err)

	if len(columns) > 0 {
		l.Printf("Column Name %v", columns[0].Name)
	}
}

func Test_main(t *testing.T) {
	main()
}
