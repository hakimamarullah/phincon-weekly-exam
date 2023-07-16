package domain

import (
	"packettrackingnet/config/consts"
	"packettrackingnet/helpers"
	"strings"
)

type Packet struct {
	PacketId    string    `json:"packetId" validate:"required"`
	Sender      Sender    `json:"sender" validate:"required"`
	Receiver    Receiver  `json:"receiver" validate:"required"`
	Origin      *Location `json:"origin" validate:"required"`
	Destination *Location `json:"destination" validate:"required"`
	Weight      float64   `json:"weight"`
	Status      string    `json:"status"`
}

func NewPacket(sender Sender, receiver Receiver, origin *Location, destination *Location, weight float64) *Packet {
	return &Packet{PacketId: helpers.GenerateUUID(), Sender: sender, Receiver: receiver, Origin: origin, Destination: destination, Weight: weight, Status: consts.PENDING}
}

type Sender struct {
	SenderId   string `json:"senderId"`
	SenderName string `json:"senderName" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
}

func NewSender(senderName string, phone string) *Sender {
	return &Sender{SenderId: helpers.GenerateUUID(), SenderName: senderName, Phone: phone}
}

type Receiver struct {
	ReceiverId   string `json:"receiverId"`
	ReceiverName string `json:"receiverName" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
}

func NewReceiver(receiverName string, phone string) *Receiver {
	return &Receiver{ReceiverId: helpers.GenerateUUID(), ReceiverName: receiverName, Phone: phone}
}

type Shipment struct {
	ShipmentId   string `json:"shipmentId"`
	*Packet      `validate:"required"`
	ShippingCost float64 `json:"shippingCost" validate:"required"`
	Service      `validate:"required"`
	CheckPoints  []*Location `json:"checkPoints"`
	IsReceived   bool        `json:"isReceived"`
}

func NewShipment(packet *Packet, shippingCost float64, service Service, checkPoints []*Location, isReceived bool) *Shipment {
	return &Shipment{ShipmentId: helpers.GenerateUUID(), Packet: packet, Service: service, ShippingCost: shippingCost, CheckPoints: checkPoints, IsReceived: isReceived}
}

type Location struct {
	LocationId   string `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address" validate:"required"`
}

func NewLocation(locationName string, address string) *Location {
	return &Location{LocationName: locationName, Address: address}
}

type Service struct {
	ServiceId        string  `json:"serviceId"`
	ServiceName      string  `json:"serviceName" validate:"required"`
	PricePerKilogram float64 `json:"pricePerKilogram" validate:"required"`
}

func NewService(serviceName string, pricePerKilogram float64) *Service {
	return &Service{ServiceId: helpers.GenerateUUID(), ServiceName: serviceName, PricePerKilogram: pricePerKilogram}
}

type PacketDetails struct {
	Packet
	IsReceived bool `json:"isReceived"`
}

func NewPacketDetails(packet Packet, isReceived bool) *PacketDetails {
	return &PacketDetails{Packet: packet, IsReceived: isReceived}
}

func (packet *Packet) GetId() string {
	return strings.ToUpper(packet.PacketId)
}

func (packet *Packet) GetOrigin() string {
	return packet.Origin.LocationName
}

func (packet *Packet) GetDestination() string {
	return packet.Destination.LocationName
}

func (packet *Packet) GetSender() string {
	return packet.Sender.SenderName
}

func (packet *Packet) GetReceiver() string {
	return packet.Receiver.ReceiverName
}

func (shipment *Shipment) GetCurrentPosition() string {
	checkPoints := shipment.CheckPoints
	if len(checkPoints) <= 0 {
		return "UNKNOWN"
	}
	return checkPoints[len(checkPoints)-1].LocationName
}

func (shipment *Shipment) GetId() string {
	return strings.ToUpper(shipment.ShipmentId)
}
