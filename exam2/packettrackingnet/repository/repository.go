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

func AddSender(sender domain.Sender) error {
	senders, err := readData[domain.Sender](config.SENDER)
	if err != nil {
		return err
	}
	*senders = append(*senders, sender)
	writeData[domain.Sender](config.SENDER, *senders)
	return nil
}

func AddReceiver(receiver domain.Receiver) error {
	receivers, err := readData[domain.Receiver](config.RECEIVER)
	if err != nil {
		return err
	}
	*receivers = append(*receivers, receiver)
	writeData[domain.Receiver](config.RECEIVER, *receivers)
	return nil
}

func AddShipment(shipment domain.Shipment) error {
	shipments, err := readData[domain.Shipment](config.SHIPMENT)
	if err != nil {
		return err
	}
	*shipments = append(*shipments, shipment)
	writeData[domain.Shipment](config.SHIPMENT, *shipments)
	return nil
}

func GetAllShipment() (*[]domain.Shipment, error) {
	return readData[domain.Shipment](config.SHIPMENT)
}

func FindShipmentById(id string) (bool, *domain.Shipment) {
	shipments, _ := readData[domain.Shipment](config.SHIPMENT)
	for i, shipment := range *shipments {
		if strings.EqualFold(shipment.Id, id) {
			return true, &(*shipments)[i]
		}
	}
	return false, nil
}

func FindSenderById(id string) (bool, *domain.Sender) {
	senders, _ := readData[domain.Sender](config.SENDER)
	for i, sender := range *senders {
		if strings.EqualFold(sender.Id, id) {
			return true, &(*senders)[i]
		}
	}
	return false, nil
}

func FindReceiverById(id string) (bool, *domain.Receiver) {
	receivers, _ := readData[domain.Receiver](config.RECEIVER)
	for i, receiver := range *receivers {
		if strings.EqualFold(receiver.Id, id) {
			return true, &(*receivers)[i]
		}
	}
	return false, nil
}

func FindPacketById(id string) (bool, *domain.Packet) {
	packets, _ := readData[domain.Packet](config.PACKET)
	for i, packet := range *packets {
		if strings.EqualFold(packet.Id, id) {
			return true, &(*packets)[i]
		}
	}
	return false, nil
}

func AddLocation(location *domain.Location) {
	locations, _ := readData[domain.Location](config.LOCATION)
	location.Id = helpers.GenerateIdLocation(len(*locations))
	*locations = append(*locations, *location)
	writeData[domain.Location](config.LOCATION, *locations)
}

func GetAllLocations() ([]domain.Location, error) {
	results, err := readData[domain.Location](config.LOCATION)
	return *results, err
}

func FindLocationById(id string) (bool, *domain.Location) {
	locations, _ := readData[domain.Location](config.LOCATION)
	for i, location := range *locations {
		if strings.EqualFold(location.Id, id) {
			return true, &(*locations)[i]
		}
	}
	return false, nil
}

func AddService(service domain.Service) error {
	services, err := readData[domain.Service](config.SERVICE)
	if err != nil {
		return err
	}
	*services = append(*services, service)
	writeData[domain.Service](config.SERVICE, *services)
	return nil
}

func GetAllServices() ([]domain.Service, error) {
	services, err := readData[domain.Service](config.SERVICE)
	return *services, err
}

func AddPacket(packet domain.Packet) error {
	packets, err := readData[domain.Packet](config.PACKET)
	if err != nil {
		return err
	}
	*packets = append(*packets, packet)
	writeData[domain.Packet](config.PACKET, *packets)
	return nil
}

func FindServiceByName(name string) (bool, *domain.Service) {
	services, _ := readData[domain.Service](config.SERVICE)
	for i, service := range *services {
		if strings.EqualFold(service.ServiceName, name) {
			return true, &(*services)[i]
		}
	}
	return false, nil
}

func FindLocationByName(name string) (bool, *domain.Location) {
	locations, _ := readData[domain.Location](config.LOCATION)
	for i, loc := range *locations {
		if strings.EqualFold(loc.LocationName, name) {
			return true, &(*locations)[i]
		}
	}
	return false, nil
}

func readData[T any](path string) (*[]T, error) {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
		return nil, errors.New("can't open file")
	}

	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var data *[]T
	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
		return errors.New("can't write data")
	}

	return nil
}
