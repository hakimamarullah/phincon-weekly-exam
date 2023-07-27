package dto

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
