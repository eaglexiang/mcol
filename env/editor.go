package env

import (
	"runtime"
)

// DefaultEditor 默认的编辑器
func DefaultEditor() (editor string) {
	switch runtime.GOOS {
	case "windows":
		editor = "notepad.exe"
	case "linux", "darwin":
		editor = "vim"
	}
	return
}
