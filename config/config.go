package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/eaglexiang/mcol/env"
	"github.com/pkg/errors"
)

// Config 配置
type Config struct {
	Addr     string `json:"addr"`
	DB       string `json:"db"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// C 全局配置
var C Config

// getFilename 获取配置文件的路径
func getFilename() (fn string, err error) {
	home, err := env.Home()
	if err != nil {
		return
	}

	fn = filepath.Join(home, ".mcol.config")

	return
}

// Load 加载配置
func Load() (err error) {
	fn, err := getFilename()
	if err != nil {
		return
	}

	buf, err := os.ReadFile(fn)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = json.Unmarshal(buf, &C)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

// Try2InitFile 初始化配置文件
func Try2InitFile(defaultConfig []byte) (filename string, err error) {
	filename, err = getFilename()
	if err != nil {
		return
	}

	ok := existConfigFile(filename)
	if ok {
		return
	}

	log.Println("初始化配置文件")
	err = os.WriteFile(filename, defaultConfig, os.FileMode(0x777))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func existConfigFile(filename string) (ok bool) {
	_, err := os.Stat(filename)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
