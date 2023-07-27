package dto

import "time"

type ResponseBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
	Code    int         `json:"code"`
}

type PacketRequest struct {
	SenderID        int64   `json:"senderId" validate:"required"`
	ReceiverID      int64   `json:"receiverId" validate:"required"`
	OriginName      string  `json:"originName" validate:"required"`
	DestinationName string  `json:"destinationName" validate:"required"`
	Weight          float64 `json:"weight" validate:"required"`
}

type ShipmentRequest struct {
	PacketID    int64  `json:"packetId" validate:"required"`
	ServiceName string `json:"serviceName" validate:"required"`
	ServiceId   int64  `json:"serviceId"`
}

type UpdateShipmentRequest struct {
	ShipmentID int64 `json:"shipmentId" validate:"required"`
	LocationID int64 `json:"locationId" validate:"required"`
}

type UpdateLocationAddressRequest struct {
	LocationName string `json:"locationName" validate:"required"`
	Address      string `json:"address" validate:"required"`
}

type ShipmentResponse struct {
	ShipmentId         int64     `json:"shipmentId"`
	PacketId           int64     `json:"packetId" validate:"required"`
	SenderId           int64     `json:"senderId"`
	SenderName         string    `json:"senderName"`
	SenderPhone        string    `json:"senderPhone"`
	ReceiverId         int64     `json:"receiverId"`
	ReceiverName       string    `json:"receiverName"`
	ReceiverPhone      string    `json:"receiverPhone"`
	OriginId           int64     `json:"originId"`
	OriginName         string    `json:"originName"`
	OriginAddress      string    `json:"originAddress"`
	DestinationId      int64     `json:"destinationId"`
	DestinationName    string    `json:"destinationName"`
	DestinationAddress string    `json:"destinationAddress" validate:"required"`
	Weight             float64   `json:"weight"`
	Status             string    `json:"status"`
	CreatedOn          time.Time `json:"createdOn"`
	UpdatedOn          time.Time `json:"updatedOn"`
}
