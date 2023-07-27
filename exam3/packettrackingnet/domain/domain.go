package domain

import (
	"packettrackingnet/config/consts"
	"time"
)

type Customer struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}
type Packet struct {
	PacketId    int64   `json:"packetId" validate:"required"`
	Sender      int64   `json:"sender" validate:"required"`
	Receiver    int64   `json:"receiver" validate:"required"`
	Origin      int64   `json:"origin" validate:"required"`
	Destination int64   `json:"destination" validate:"required"`
	Weight      float64 `json:"weight"`
	Status      string  `json:"status"`
}

func NewPacket(sender int64, receiver int64, origin int64, destination int64, weight float64) *Packet {
	return &Packet{Sender: sender, Receiver: receiver, Origin: origin, Destination: destination, Weight: weight, Status: consts.PENDING}
}

func NewCustomer(senderName string, phone string) *Customer {
	return &Customer{Name: senderName, Phone: phone}
}

type Shipment struct {
	ShipmentId   int64     `json:"shipmentId"`
	Packet       int64     `validate:"required"`
	ShippingCost float64   `json:"shippingCost" validate:"required"`
	Service      int64     `validate:"required"`
	CheckPoints  []int64   `json:"checkPoints"`
	IsReceived   bool      `json:"isReceived"`
	CreatedOn    time.Time `json:"createdOn"`
	UpdatedOn    time.Time `json:"updatedOn"`
}

func NewShipment(packet int64, shippingCost float64, service int64, checkPoints []int64, isReceived bool) *Shipment {
	return &Shipment{Packet: packet, Service: service, ShippingCost: shippingCost, CheckPoints: checkPoints, IsReceived: isReceived}
}

type Location struct {
	LocationId   int64  `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address" validate:"required"`
}

func NewLocation(locationName string, address string) *Location {
	return &Location{LocationName: locationName, Address: address}
}

type Service struct {
	ServiceId        int64   `json:"serviceId"`
	ServiceName      string  `json:"serviceName" validate:"required"`
	PricePerKilogram float64 `json:"pricePerKilogram" validate:"required"`
}

func NewService(serviceName string, pricePerKilogram float64) *Service {
	return &Service{ServiceName: serviceName, PricePerKilogram: pricePerKilogram}
}

func (packet *Packet) GetId() int64 {
	return packet.PacketId
}

func (packet *Packet) GetOrigin() int64 {
	return packet.Origin
}

func (packet *Packet) GetDestination() int64 {
	return packet.Destination
}

func (packet *Packet) GetSender() int64 {
	return packet.Sender
}

func (packet *Packet) GetReceiver() int64 {
	return packet.Receiver
}

func (shipment *Shipment) GetCurrentPosition() int64 {
	checkPoints := shipment.CheckPoints
	if len(checkPoints) <= 0 {
		return -1
	}
	return checkPoints[len(checkPoints)-1]
}

func (shipment *Shipment) GetId() int64 {
	return shipment.ShipmentId
}
