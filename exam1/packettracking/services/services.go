package services

import (
	"errors"
	"packettracking/domain"
	"packettracking/repository"
	"strings"
)

func CreateShipment(packet domain.Packet, service domain.Service) string {
	shippingCost := packet.Weight * service.PricePerKilogram
	newShipment := domain.NewShipment(packet, shippingCost, service, []domain.Location{}, false)
	senderShipments := packet.Sender.Shipments
	repository.AddShipment(*newShipment)
	senderShipments = append(senderShipments, *newShipment)
	return newShipment.Id
}

func GetAllReceivedPackets() []domain.Packet {
	shipments := repository.GetAllShipment()
	var results []domain.Packet

	for _, item := range shipments {
		if item.IsReceived {
			results = append(results, item.Packet)
		}
	}

	return results
}

func GetAllPacketsByLocationName(locationName string) []domain.PacketDetails {
	shipments := repository.GetAllShipment()
	var results []domain.PacketDetails

	for _, ship := range shipments {
		for _, loc := range ship.CheckPoints {
			if strings.EqualFold(loc.LocationName, locationName) {
				results = append(results, *domain.NewPacketDetails(ship.Packet, ship.IsReceived))
			}
		}
	}

	return results
}
func GetAllCheckpoints() []domain.Location {
	return repository.GetAllLocations()
}

func UpdateShipmentCheckpoint(shipmentId string, locationId string) (*domain.Shipment, error) {
	exist, shipment := repository.FindShipmentById(shipmentId)
	existLocation, loc := repository.FindLocationById(locationId)
	if !exist || !existLocation {
		return nil, errors.New("data not found")
	}

	shipment.CheckPoints = append(shipment.CheckPoints, *loc)
	for _, item := range shipment.CheckPoints {
		if item.Id == shipment.Packet.Destination.Id {
			shipment.IsReceived = true
		}
	}

	return shipment, nil
}

func GetAllShipment() []domain.Shipment {
	return repository.GetAllShipment()
}

func CreateService(serviceName string, pricePerKilogram float64) {
	repository.AddService(*domain.NewService(serviceName, pricePerKilogram))
}

func AddSender(senderName string, phone string) *domain.Sender {
	newSender := domain.NewSender(senderName, phone)
	repository.AddSender(*newSender)
	return newSender
}

func AddReceiver(receiverName string, phone string) *domain.Receiver {
	newReceiver := domain.NewReceiver(receiverName, phone)
	repository.AddReceiver(*newReceiver)
	return newReceiver
}

func AddLocation(locationName string, address string) *domain.Location {
	exist, loc := repository.FindLocationByName(locationName)
	if !exist {
		newLocation := domain.NewLocation(locationName, address)
		repository.AddLocation(newLocation)
		return newLocation
	}
	return loc
}

func AddPacket(sender domain.Sender, receiver domain.Receiver, destination domain.Location, weight float64) *domain.Packet {
	newPacket := domain.NewPacket(sender, receiver, destination, weight)
	repository.AddPacket(*newPacket)
	return newPacket
}

func FindPacketById(id string) *domain.Packet {
	exist, packet := repository.FindPacketById(id)
	if !exist {
		return nil
	}

	return packet
}

func GetAllServices() []domain.Service {
	return repository.GetAllServices()
}

func GetAllServiceNames() []string {
	var results []string

	for _, item := range repository.GetAllServices() {
		results = append(results, item.ServiceName)
	}
	return results
}

func GetServiceByName(name string) *domain.Service {
	exist, service := repository.FindServiceByName(name)
	if !exist {
		return nil
	}

	return service
}

func GetLocationByName(name string) *domain.Location {
	exist, location := repository.FindLocationByName(name)
	if !exist {
		return nil
	}
	return location
}

func GetShipmentById(id string) (bool, *domain.Shipment) {
	return repository.FindShipmentById(id)
}
