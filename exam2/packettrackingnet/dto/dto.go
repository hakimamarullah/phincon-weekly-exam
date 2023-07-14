package dto

type ResponseBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count"`
	Code    int         `json:"code"`
}

type PacketRequest struct {
	SenderID        string  `json:"senderId" validate:"required"`
	ReceiverID      string  `json:"receiverId" validate:"required"`
	OriginName      string  `json:"originName" validate:"required"`
	DestinationName string  `json:"destinationName" validate:"required"`
	Weight          float64 `json:"weight" validate:"required"`
}

type ShipmentRequest struct {
	PacketID    string `json:"packetId" validate:"required"`
	ServiceName string `json:"serviceName" validate:"required"`
}

type UpdateShipmentRequest struct {
	ShipmentID string `json:"shipmentId" validate:"required"`
	LocationID string `json:"locationId" validate:"required"`
}

type UpdateLocationAddressRequest struct {
	LocationName string `json:"locationName" validate:"required"`
	Address      string `json:"address" validate:"required"`
}
