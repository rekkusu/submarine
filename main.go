package main

import (
	"github.com/activedefense/submarine/adctf"
	"github.com/activedefense/submarine/adctf/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server := adctf.New()
	server.Start(config.App.Listen)
}
