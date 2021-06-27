package config

import (
	"encoding/json"
	"os"

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

// Load 加载配置
func Load() (err error) {
	buf, err := os.ReadFile("mcol.config")
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
