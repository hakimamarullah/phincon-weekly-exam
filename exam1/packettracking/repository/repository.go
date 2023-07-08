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

func GetAllSenders() []domain.Sender {
	return senders
}

func FindSenderById(id string) (bool, *domain.Sender) {
	for _, sender := range senders {
		if sender.Id == id {
			return true, &sender
		}
	}
	return false, nil
}

func AddReceiver(receiver domain.Receiver) {
	receivers = append(receivers, receiver)
}

func GetAllReceiver() []domain.Receiver {
	return receivers
}

func FindReceiverById(id string) (bool, *domain.Receiver) {
	for _, receiver := range receivers {
		if receiver.Id == id {
			return true, &receiver
		}
	}
	return false, nil
}

func AddShipment(shipment domain.Shipment) {
	shipments = append(shipments, shipment)
}

func GetAllShipment() []domain.Shipment {
	return shipments
}

func FindShipmentById(id string) (bool, *domain.Shipment) {
	for _, shipment := range shipments {
		if strings.EqualFold(shipment.Id, id) {
			return true, &shipment
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
	for _, location := range locations {
		if strings.EqualFold(location.Id, id) {
			return true, &location
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

func FindServiceById(id string) (bool, *domain.Service) {
	for _, service := range services {
		if service.Id == id {
			return true, &service
		}
	}
	return false, nil
}

func AddPacket(packet domain.Packet) {
	packets = append(packets, packet)
}

func GetAllPackets() []domain.Packet {
	return packets
}

func FindPacketById(id string) (bool, *domain.Packet) {
	for _, packet := range packets {
		if packet.Id == id {
			return true, &packet
		}
	}
	return false, nil
}

func FindServiceByName(name string) (bool, *domain.Service) {
	for _, service := range services {
		if strings.EqualFold(service.ServiceName, name) {
			return true, &service
		}
	}
	return false, nil
}

func FindLocationByName(name string) (bool, *domain.Location) {
	for _, loc := range locations {
		if strings.EqualFold(loc.LocationName, name) {
			return true, &loc
		}
	}
	return false, nil
}
