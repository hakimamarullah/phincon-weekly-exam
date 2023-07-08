package domain

import (
	"packettracking/utils"
)

type Packet struct {
	Id          string
	Sender      Sender
	Receiver    Receiver
	Origin      Location
	Destination Location
	Weight      float64
}

func NewPacket(sender Sender, receiver Receiver, origin Location, destination Location, weight float64) *Packet {
	return &Packet{Id: utils.GenerateId(), Sender: sender, Receiver: receiver, Origin: origin, Destination: destination, Weight: weight}
}

type Sender struct {
	Id         string
	SenderName string
	Phone      string
}

func NewSender(senderName string, phone string) *Sender {
	return &Sender{Id: utils.GenerateId(), SenderName: senderName, Phone: phone}
}

type Receiver struct {
	Id           string
	ReceiverName string
	Phone        string
}

func NewReceiver(receiverName string, phone string) *Receiver {
	return &Receiver{Id: utils.GenerateId(), ReceiverName: receiverName, Phone: phone}
}

type Shipment struct {
	Id           string
	Packet       Packet
	ShippingCost float64
	Service      Service
	CheckPoints  []Location
	IsReceived   bool
}

func NewShipment(packet Packet, shippingCost float64, service Service, checkPoints []Location, isReceived bool) *Shipment {
	return &Shipment{Id: utils.GenerateId(), Packet: packet, Service: service, ShippingCost: shippingCost, CheckPoints: checkPoints, IsReceived: isReceived}
}

type Location struct {
	Id           string
	LocationName string
	Address      string
}

func NewLocation(locationName string, address string) *Location {
	return &Location{LocationName: locationName, Address: address}
}

type Service struct {
	Id               string
	ServiceName      string
	PricePerKilogram float64
}

func NewService(serviceName string, pricePerKilogram float64) *Service {
	return &Service{Id: utils.GenerateId(), ServiceName: serviceName, PricePerKilogram: pricePerKilogram}
}

type PacketDetails struct {
	Packet     Packet
	IsReceived bool
}

func NewPacketDetails(packet Packet, isReceived bool) *PacketDetails {
	return &PacketDetails{Packet: packet, IsReceived: isReceived}
}
