package main

import (
	"github.com/activedefense/submarine/adctf"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server := adctf.New(adctf.ADCTFConfig{
		DriverName:     "sqlite3",
		DataSourceName: ":memory:?parseTime=true",
		JWTSecret:      []byte("e81061ace8c9e0b568c00075ecda0d8c42d"),
		Debug:          true,
	})
	server.Start("127.0.0.1:8000")
}
