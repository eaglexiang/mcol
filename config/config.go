package config

import (
	"encoding/json"
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
