package config

import "strings"

var whiteList = []string{"/users", "/login"}

func WhiteListed(path string) bool {
	for _, item := range whiteList {
		if strings.EqualFold(item, path) {
			return true
		}
	}
	return false
}
