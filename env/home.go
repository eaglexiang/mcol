package env

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

// ErrHomeNotFound HOME 目录找不到
var ErrHomeNotFound = errors.New("home path not found")

// Home 获取用户 Home 目录
func Home() (homePath string, err error) {
	user, err := user.Current()
	if err == nil {
		homePath = user.HomeDir
		return
	}

	homePath, err = unixHome()
	if err == nil {
		return
	}

	homePath, err = windowsHome()
	return
}

func unixHome() (homePath string, err error) {
	homePath = os.Getenv("HOME")
	if homePath != "" {
		return
	}

	homePath = os.Getenv("~")
	if homePath != "" {
		return
	}

	err = errors.WithStack(ErrHomeNotFound)
	return
}

func windowsHome() (homePath string, err error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	if drive != "" && path != "" {
		homePath = filepath.Join(drive, path)
		return
	}

	homePath = os.Getenv("USERPROFILE")
	if homePath == "" {
		err = errors.WithStack(ErrHomeNotFound)
	}
	return
}
