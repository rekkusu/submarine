package main

import (
	"github.com/activedefense/submarine/adctf"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server := adctf.New()
	server.Start("127.0.0.1:8000")
}
