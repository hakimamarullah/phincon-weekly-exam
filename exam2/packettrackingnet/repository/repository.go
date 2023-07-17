package repository

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"packettrackingnet/config"
	"packettrackingnet/domain"
	"packettrackingnet/helpers"
	"strings"
)

var senders = make([]domain.Sender, 0)
var receivers = make([]domain.Receiver, 0)
var shipments = make([]domain.Shipment, 0)
var locations = make([]domain.Location, 0)
var services = make([]domain.Service, 0)
var packets = make([]domain.Packet, 0)

func InitDatastore() {
	dataSender, _ := readData[domain.Sender](config.SENDER)
	senders = *dataSender

	dataReceiver, _ := readData[domain.Receiver](config.RECEIVER)
	receivers = *dataReceiver

	dataShipments, _ := readData[domain.Shipment](config.SHIPMENT)
	shipments = *dataShipments

	dataLocation, _ := readData[domain.Location](config.LOCATION)
	locations = *dataLocation

	dataService, _ := readData[domain.Service](config.SERVICE)
	services = *dataService

	dataPackets, _ := readData[domain.Packet](config.PACKET)
	packets = *dataPackets
}

func PersistData() {
	writeData[domain.Sender](config.SENDER, senders)
	writeData[domain.Receiver](config.RECEIVER, receivers)
	writeData[domain.Shipment](config.SHIPMENT, shipments)
	writeData[domain.Location](config.LOCATION, locations)
	writeData[domain.Service](config.SERVICE, services)
	writeData[domain.Packet](config.PACKET, packets)
}

func AddSender(sender domain.Sender) {
	senders = append(senders, sender)
}

func AddReceiver(receiver domain.Receiver) {
	receivers = append(receivers, receiver)
}

func AddShipment(shipment domain.Shipment) {
	shipments = append(shipments, shipment)
}

func GetAllShipment() *[]domain.Shipment {
	return &shipments
}

func FindShipmentById(id string) (bool, *domain.Shipment) {
	for i, shipment := range shipments {
		if strings.EqualFold(shipment.ShipmentId, id) {
			return true, &shipments[i]
		}
	}
	return false, nil
}

func FindSenderById(id string) (bool, *domain.Sender) {
	for i, sender := range senders {
		if strings.EqualFold(sender.SenderId, id) {
			return true, &senders[i]
		}
	}
	return false, nil
}

func FindReceiverById(id string) (bool, *domain.Receiver) {
	for i, receiver := range receivers {
		if strings.EqualFold(receiver.ReceiverId, id) {
			return true, &receivers[i]
		}
	}
	return false, nil
}

func FindPacketById(id string) (bool, *domain.Packet) {
	for i, packet := range packets {
		if strings.EqualFold(packet.PacketId, id) {
			return true, &packets[i]
		}
	}
	return false, nil
}

func AddLocation(location *domain.Location) {
	location.LocationId = helpers.GenerateIdLocation(len(locations))
	locations = append(locations, *location)
}

func GetAllLocations() ([]domain.Location, error) {
	return locations, nil
}

func FindLocationById(id string) (bool, *domain.Location) {
	for i, location := range locations {
		if strings.EqualFold(location.LocationId, id) {
			return true, &locations[i]
		}
	}
	return false, nil
}

func AddService(service domain.Service) {
	services = append(services, service)
}

func GetAllServices() ([]domain.Service, error) {
	return services, nil
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

func readData[T any](path string) (*[]T, error) {
	reader, err := os.Open(path)
	if err != nil {
		helpers.LogError(err)
		return nil, errors.New("can't open file")
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var data *[]T
	if err := decoder.Decode(&data); err != nil {
		helpers.LogError(err)
		return nil, errors.New("error reading data")
	}

	return data, nil
}

func writeData[T any](path string, data []T) error {
	writer, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	defer writer.Close()

	encoder := json.NewEncoder(writer)

	if err := encoder.Encode(data); err != nil {
		helpers.LogError(err)
		return errors.New("can't write data")
	}

	return nil
}

func ExistPacketShipment(packetId string) bool {
	for _, item := range shipments {
		if strings.EqualFold(item.PacketId, packetId) {
			return true
		}
	}
	return false
}
