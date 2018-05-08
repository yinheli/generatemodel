# generate model

automatic generate go struct from database(MySQL) table

## usage

```
go get github.com/yinheli/generatemodel
```

add comment in `model` package

eg: `model/db.go`
```
//go:generate generatemodel
package model

// ... some code ...
```

```
export database_uri="user:pass@(host:3306)/db_name?charset=utf8&parseTime=True&loc=Local
go generate ./model
```