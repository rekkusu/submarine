package main

import (
	"flag"
	"os"
	"strings"

	"github.com/activedefense/submarine/adctf"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var driver string
	var source string
	var secret string
	var listen string
	var debug bool
	flag.StringVar(&driver, "driver", "sqlite3", "DB Driver Name")
	flag.StringVar(&source, "source", "submarine.db?parseTime=true", "DB Source Name")
	flag.StringVar(&secret, "secret", "e81061ace8c9e0b568c00075ecda0d8c42d", "Application SecretKey")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&listen, "listen", "127.0.0.1:8000", "Host/port to listen")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper("submarine_" + f.Name)); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	server := adctf.New(adctf.ADCTFConfig{
		DriverName:     driver,
		DataSourceName: source,
		JWTSecret:      []byte(secret),
		Debug:          debug,
	})
	server.Start(listen)
}
