package mapper

import (
	"packettrackingnet/dao"
	"packettrackingnet/dto"
)

func ShipmentDaoToShipmentResponse(dao dao.ShipmentDAO) *dto.ShipmentResponse {
	return &dto.ShipmentResponse{
		ShipmentId:         dao.ShipmentId,
		PacketId:           dao.Packet.PacketId,
		SenderId:           dao.Packet.Sender.CustomerId,
		SenderName:         dao.Packet.Sender.Name,
		SenderPhone:        dao.Packet.Sender.Phone,
		ReceiverId:         dao.Packet.Receiver.CustomerId,
		ReceiverName:       dao.Packet.Receiver.Name,
		ReceiverPhone:      dao.Packet.Receiver.Phone,
		OriginId:           dao.Packet.Origin.LocationId,
		OriginName:         dao.Packet.Origin.LocationName,
		OriginAddress:      dao.Packet.Origin.Address,
		DestinationId:      dao.Packet.Destination.LocationId,
		DestinationName:    dao.Packet.Destination.LocationName,
		DestinationAddress: dao.Packet.Destination.Address,
		Weight:             dao.Packet.Weight,
		Status:             dao.Packet.Status,
		CreatedOn:          dao.CreatedOn,
		UpdatedOn:          dao.UpdatedOn,
	}
}
