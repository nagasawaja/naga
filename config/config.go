package config

import (
	"github.com/jinzhu/configor"
	"os"
	"path/filepath"
	"time"
)

var Config = &config{}

type config struct {
	AppName    string `yaml:"appName"`
	LogPath    string `yaml:"logPath"`
	ServerAddr string `yaml:"serverAddr"`

	Mysql struct {
		Host      string
		User      string
		Password  string
		Port      string
		Database  string
		Charset   string
		MaxActive int
		MaxIdle   int
		Logdebug  bool
	}
	Redis struct {
		Host        string
		Password    string
		Port        string
		MaxActive   int
		MaxIdle     int
		Idletimeout time.Duration
		Database    int
	}
}

func (c *config) Start(fileName string) {
	var a = configor.New(&configor.Config{})
	err := a.Load(c, fileName)
	if err != nil {
		panic(err)
	}
	if c.LogPath == "" {
		dir, _ := os.Executable()
		exPath := filepath.Dir(dir)
		c.LogPath = filepath.FromSlash(exPath + "/logPath/")
	}
}
