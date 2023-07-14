package domain

import (
	"packettrackingnet/helpers"
)

type Packet struct {
	Id          string   `json:"id"`
	Sender      Sender   `json:"sender" validate:"required"`
	Receiver    Receiver `json:"receiver" validate:"required"`
	Origin      Location `json:"origin" validate:"required"`
	Destination Location `json:"destination" validate:"required"`
	Weight      float64  `json:"weight"`
}

func NewPacket(sender Sender, receiver Receiver, origin Location, destination Location, weight float64) *Packet {
	return &Packet{Id: helpers.GenerateUUID(), Sender: sender, Receiver: receiver, Origin: origin, Destination: destination, Weight: weight}
}

type Sender struct {
	Id         string `json:"id"`
	SenderName string `json:"senderName" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
}

func NewSender(senderName string, phone string) *Sender {
	return &Sender{Id: helpers.GenerateUUID(), SenderName: senderName, Phone: phone}
}

type Receiver struct {
	Id           string `json:"id"`
	ReceiverName string `json:"receiverName" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
}

func NewReceiver(receiverName string, phone string) *Receiver {
	return &Receiver{Id: helpers.GenerateUUID(), ReceiverName: receiverName, Phone: phone}
}

type Shipment struct {
	Id           string     `json:"id"`
	Packet       Packet     `json:"packet" validate:"required"`
	ShippingCost float64    `json:"shippingCost" validate:"required"`
	Service      Service    `json:"service" validate:"required"`
	CheckPoints  []Location `json:"checkPoints"`
	IsReceived   bool       `json:"isReceived"`
}

func NewShipment(packet Packet, shippingCost float64, service Service, checkPoints []Location, isReceived bool) *Shipment {
	return &Shipment{Id: helpers.GenerateUUID(), Packet: packet, Service: service, ShippingCost: shippingCost, CheckPoints: checkPoints, IsReceived: isReceived}
}

type Location struct {
	Id           string `json:"id"`
	LocationName string `json:"locationName"`
	Address      string `json:"address" validate:"required"`
}

func NewLocation(locationName string, address string) *Location {
	return &Location{LocationName: locationName, Address: address}
}

type Service struct {
	Id               string  `json:"id"`
	ServiceName      string  `json:"serviceName" validate:"required"`
	PricePerKilogram float64 `json:"pricePerKilogram" validate:"required"`
}

func NewService(serviceName string, pricePerKilogram float64) *Service {
	return &Service{Id: helpers.GenerateUUID(), ServiceName: serviceName, PricePerKilogram: pricePerKilogram}
}

type PacketDetails struct {
	Packet     Packet
	IsReceived bool
}

func NewPacketDetails(packet Packet, isReceived bool) *PacketDetails {
	return &PacketDetails{Packet: packet, IsReceived: isReceived}
}
