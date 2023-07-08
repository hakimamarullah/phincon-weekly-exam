package services

import (
	"errors"
	"packettracking/domain"
	"packettracking/repository"
	"strings"
)

// CreateShipment function to create shipment and save it to the repository.
// Returns the id of newly created shipment.
func CreateShipment(packet domain.Packet, service domain.Service) string {
	shippingCost := packet.Weight * service.PricePerKilogram
	newShipment := domain.NewShipment(packet, shippingCost, service, []domain.Location{packet.Origin}, false)
	repository.AddShipment(*newShipment)
	return newShipment.Id
}

// GetAllReceivedPackets function to get all packets that has been delivered.
// It checks the status of the Shipment-received flag.
// Returns slice of Packet
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

// GetAllPacketsByLocationName function to get all the packets that has been going through
// specified location name.
// It returns a slice of PacketDetails struct.
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

// GetAllCheckpoints function to get all Location saved in repository.
// Returns slice of Location.
func GetAllCheckpoints() []domain.Location {
	return repository.GetAllLocations()
}

// UpdateShipmentCheckpoint function to update checkpoint of shipment.
// It only appends the new checkpoint to the end of the checkpoint's slice.
// Returns error if the Shipment or Location is not found.
func UpdateShipmentCheckpoint(shipmentId string, locationId string) (*domain.Shipment, error) {
	exist, shipment := repository.FindShipmentById(shipmentId)
	existLocation, loc := repository.FindLocationById(locationId)
	if !exist || !existLocation {
		return nil, errors.New("data not found")
	}

	shipment.CheckPoints = append(shipment.CheckPoints, *loc)

	if len(shipment.CheckPoints) > 0 {
		lastPos := shipment.CheckPoints[len(shipment.CheckPoints)-1]
		if strings.EqualFold(lastPos.Id, shipment.Packet.Destination.Id) {
			shipment.IsReceived = true
		}
	}
	return shipment, nil
}

// GetAllShipment function to get all shipments saved in repository.
// Returns slice of Shipment struct.
func GetAllShipment() []domain.Shipment {
	return repository.GetAllShipment()
}

// CreateService procedure to create and save a new Service into repository.
func CreateService(serviceName string, pricePerKilogram float64) {
	repository.AddService(*domain.NewService(serviceName, pricePerKilogram))
}

// AddSender function to add new sender into repository.
// Returns the newly created Sender pointer.
func AddSender(senderName string, phone string) *domain.Sender {
	newSender := domain.NewSender(senderName, phone)
	repository.AddSender(*newSender)
	return newSender
}

// AddReceiver function to add new receiver into the repository.
// Returns the newly created Receiver pointer.
func AddReceiver(receiverName string, phone string) *domain.Receiver {
	newReceiver := domain.NewReceiver(receiverName, phone)
	repository.AddReceiver(*newReceiver)
	return newReceiver
}

// AddLocation function to add new Location into the repository.
// Returns the newly created Location pointer.
// It only puts the Location into the repository if the new Location name doesn't exist in the repository.
func AddLocation(locationName string, address string) (bool, *domain.Location) {
	exist, loc := repository.FindLocationByName(locationName)
	if !exist {
		newLocation := domain.NewLocation(locationName, address)
		repository.AddLocation(newLocation)
		return true, newLocation
	}
	return false, loc
}

// AddPacket function to add the new packet into the repository.
// Returns the newly created Packet pointer.
func AddPacket(sender domain.Sender, receiver domain.Receiver, origin domain.Location, destination domain.Location, weight float64) *domain.Packet {
	newPacket := domain.NewPacket(sender, receiver, origin, destination, weight)
	repository.AddPacket(*newPacket)
	return newPacket
}

// GetAllServices function to get all Services from the repository.
// Returns slice of Services.
func GetAllServices() []domain.Service {
	return repository.GetAllServices()
}

// GetAllServiceNames function to get all service names from repository,
// Returns slice of service's name
func GetAllServiceNames() []string {
	var results []string
	services := repository.GetAllServices()
	for _, item := range services {
		results = append(results, item.ServiceName)
	}
	return results
}

// GetServiceByName function to query service from repository by service name.
// This method is case-insensitive.
// Returns a single Service pointer.
func GetServiceByName(name string) *domain.Service {
	_, service := repository.FindServiceByName(name)
	return service
}

// GetShipmentById function to query Shipment from repository by Shipment ID.
// Returns boolean to indicate the existence of Shipment and Shipment pointer.
func GetShipmentById(id string) (bool, *domain.Shipment) {
	return repository.FindShipmentById(id)
}
