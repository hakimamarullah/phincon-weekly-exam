package repository

import (
	"packettracking/domain"
	"packettracking/utils"
	"strings"
)

var senders []domain.Sender
var receivers []domain.Receiver
var shipments []domain.Shipment
var locations []domain.Location
var services []domain.Service
var packets []domain.Packet

func AddSender(sender domain.Sender) {
	senders = append(senders, sender)
}

func AddReceiver(receiver domain.Receiver) {
	receivers = append(receivers, receiver)
}

func AddShipment(shipment domain.Shipment) {
	shipments = append(shipments, shipment)
}

func GetAllShipment() []domain.Shipment {
	return shipments
}

func FindShipmentById(id string) (bool, *domain.Shipment) {
	for i, shipment := range shipments {
		if strings.EqualFold(shipment.Id, id) {
			return true, &shipments[i]
		}
	}
	return false, nil
}

func AddLocation(location *domain.Location) {
	location.Id = utils.GenerateIdLocation(len(locations))
	locations = append(locations, *location)
}

func GetAllLocations() []domain.Location {
	return locations
}

func FindLocationById(id string) (bool, *domain.Location) {
	for i, location := range locations {
		if strings.EqualFold(location.Id, id) {
			return true, &locations[i]
		}
	}
	return false, nil
}

func AddService(service domain.Service) {
	services = append(services, service)
}

func GetAllServices() []domain.Service {
	return services
}

func AddPacket(packet domain.Packet) {
	packets = append(packets, packet)
}

func FindServiceByName(name string) (bool, *domain.Service) {
	for i, service := range services {
		if strings.EqualFold(service.ServiceName, name) {
			return true, &services[i]
		}
	}
	return false, nil
}

func FindLocationByName(name string) (bool, *domain.Location) {
	for i, loc := range locations {
		if strings.EqualFold(loc.LocationName, name) {
			return true, &locations[i]
		}
	}
	return false, nil
}
