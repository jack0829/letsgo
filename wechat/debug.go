package wechat

import (
	"os"
	"strings"
)

var debug bool

func Debug() bool {
	return debug
}

func init() {
	switch strings.ToLower(os.Getenv("WECHAT_DEBUG")) {
	case "true", "yes", "y", "on", "enable", "1":
		debug = true
	}
}
