package main

import (
	"flag"
	"os"
	"strings"

	"github.com/activedefense/submarine/adctf"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var config adctf.ADCTFConfig
	var secret string
	var listen string
	flag.StringVar(&config.DriverName, "driver", "sqlite3", "DB Driver Name")
	flag.StringVar(&config.DataSourceName, "source", "submarine.db?parseTime=true", "DB Source Name")
	flag.StringVar(&secret, "secret", "e81061ace8c9e0b568c00075ecda0d8c42d", "Application SecretKey")
	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")
	flag.StringVar(&listen, "listen", "127.0.0.1:8000", "Host/port to listen")
	flag.StringVar(&config.MasterPassword, "password", "masterpassword", "Master Password")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper("submarine_" + f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()
	config.JWTSecret = []byte(secret)

	server := adctf.New(config)
	server.Start(listen)
}
