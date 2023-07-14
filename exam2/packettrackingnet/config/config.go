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

const (
	LOCATION = "./repository/location.json"
	PACKET   = "./repository/packet.json"
	RECEIVER = "./repository/receiver.json"
	SENDER   = "./repository/sender.json"
	SERVICE  = "./repository/service.json"
	SHIPMENT = "./repository/shipment.json"
)
