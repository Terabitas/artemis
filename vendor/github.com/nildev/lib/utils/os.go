package utils

import (
	"os"
	"os/user"
	"strings"
)

// ResolveHome resolves HomeDir
func ResolveHome(path string) string {
	if path[:2] == "~/" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = strings.Replace(path, "~/", dir, 1)
	}

	return path
}

// Exists method
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
