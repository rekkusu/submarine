package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "submarine.toml", "Specify config file")
	flag.Parse()

	viper.SetDefault("db.driver", "sqlite3")
	viper.SetDefault("db.source", "submarine.db?parseTime=true")
	viper.SetDefault("app.secret", "e81061ace8c9e0b568c00075ecda0d8c42d")
	viper.SetDefault("app.admin_password", "masterpassword")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("app.listen", "127.0.0.1:8000")
	viper.SetDefault("ctf.team", true)
	viper.AddConfigPath("./")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	viper.Sub("db").Unmarshal(&DB)
	viper.Sub("app").Unmarshal(&App)
	viper.Sub("ctf").Unmarshal(&CTF)
}

var DB struct {
	Driver string
	Source string
}
var App struct {
	Secret        []byte
	AdminPassword string
	Debug         bool
	Listen        string
}
var CTF struct {
	Team bool
}
