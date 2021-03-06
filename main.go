package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	l = log.New(os.Stdout, "[generatemodel] ", log.LstdFlags)

	typeMap     = flag.String("typeMap", "", "mapping std type to alias (type)")
	typeMapping = make(map[string]string)

	db *sql.DB
)

func main() {
	gofile := os.Getenv("GOFILE")
	if gofile == "" {
		l.Fatal("GOFILE env not find")
	}

	dir := path.Dir(gofile)

	// overwrite
	overwrite := os.Getenv("gm_overwrite") == "true"

	flag.Parse()
	if !flag.Parsed() {
		l.Fatal("flag parse faild")
	}

	if *typeMap != "" {
		validTypeMap := regexp.MustCompile(`.*:.*(,.*:.*)*`)
		if !validTypeMap.MatchString(*typeMap) {
			l.Fatal("typeMap param format error")
		}

		for _, v := range strings.Split(*typeMap, ",") {
			item := strings.Split(v, ":")
			typeMapping[item[0]] = item[1]
		}
	}

	if err := openDB(); err != nil {
		l.Fatal(err)
	}
	defer db.Close()

	tables, err := tables()
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {

		fileName := fmt.Sprintf("%s.go", strings.ToLower(table.Name))
		file := path.Join(dir, fileName)
		if _, err := os.Stat(file); err == nil && !overwrite {
			l.Printf("file %s exists ignore generate", fileName)
			continue
		}

		table.Columns, err = columns(table.Name)
		if err != nil {
			log.Fatal(err)
		}

		l.Printf("process table: %s", table.Name)
		data, err := table.ToStruct()
		if err != nil {
			l.Fatal(err)
		}

		err = ioutil.WriteFile(file, data, 0600)
		if err != nil {
			l.Fatal(err)
		}
	}
}

func openDB() error {
	dsn := os.Getenv("database_uri")

	if dsn == "" {
		return errors.New("dsn is empty")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func tables() ([]*Table, error) {
	rows, err := db.Query("select t.table_name Name, t.table_comment Comment from information_schema.tables t where t.table_schema = database()")
	if err != nil {
		return nil, err
	}

	var tables []*Table
	for rows.Next() {
		var t Table
		err = rows.Scan(&t.Name, &t.Comment)
		if err != nil {
			return nil, err
		}
		t.TitleCaseName = TitleCase(t.Name)
		tables = append(tables, &t)
	}

	return tables, err
}

func columns(table string) ([]*Column, error) {
	rows, err := db.Query("select column_name Name, column_type, column_comment Comment, lower(is_nullable) is_nullable from information_schema.Columns t where t.table_schema=database() and t.table_name=?", table)

	var cols []*Column
	for rows.Next() {
		var (
			c        Column
			nullable string
		)
		err = rows.Scan(&c.Name, &c.DataType, &c.Comment, &nullable)
		if err != nil {
			return nil, err
		}
		c.Nullable = nullable == "yes"
		c.GoType = DataType(c.DataType, c.Nullable, typeMapping)
		c.TitleCaseName = TitleCase(c.Name)
		c.CamelCaseName = CamelCase(c.Name)
		c.Tag = Tag(c)
		cols = append(cols, &c)
	}
	return cols, nil
}
