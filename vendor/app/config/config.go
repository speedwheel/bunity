package config

import (
	"github.com/BurntSushi/toml"
	"runtime"
    "os"
    "path/filepath"
)

type AppConfig struct {
	Title string
	DB database `toml:"database"`
	User user `toml:"user"`
}

type database struct {
	Mongo mongo
} 

type mongo struct {
	Username string
	Password string
	Server string
	Database string
}

type user struct {
	SecretKey string `toml:"secretkey"`
}

var mainConfig *AppConfig
const packagePath = "/src/github.com/speedwheel/kazeli/"

func Init() *AppConfig {
	if mainConfig == nil {
		appPath := GetAppPath()
		if _, err := toml.DecodeFile(appPath+"config/config.toml", &mainConfig); err != nil {
			panic(err.Error())
		}
	}
	return mainConfig
}

func GetAppPath() string {
	return defaultGOPATH()+packagePath
}
	

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			return ""
		}
		return def
	}
	return ""
}