package repository

import (
	"packettrackingnet/config/consts"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
)

func AddPacket(packet domain.Packet) (int64, error) {
	result, err := DB.Exec("INSERT INTO Packet(Sender, Receiver, Origin, Destination, Weight, Status) VALUES(?,?,?,?,?,?)", packet.Sender, packet.Receiver, packet.GetOrigin(), packet.GetDestination(), packet.Weight, packet.Status)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func UpdatePacketStatus(packetId int64, status string) error {
	_, err := DB.Exec("UPDATE Packet SET Status = ? WHERE Id = ?", status, packetId)
	if err != nil {
		return err
	}
	return nil
}

func GetAllReceivedPackets() ([]dao.PacketDAO, error) {
	var packets = make([]dao.PacketDAO, 0)
	query := `
		SELECT 
			p.Id AS PacketId,
			p.Weight,
			p.Status,
			cSender.CustomerId AS SenderId,
			cSender.Name AS SenderName,
			cSender.Phone AS SenderPhone,
			cReceiver.CustomerId AS ReceiverId,
			cReceiver.Name AS ReceiverName,
			cReceiver.Phone AS ReceiverPhone,
			lOrigin.LocationId AS OriginId,
			lOrigin.Name AS OriginName,
			lOrigin.Address AS OriginAddress,
			lDestination.LocationId AS DestinationId,
			lDestination.Name AS DestinationName,
			lDestination.Address AS DestinationAddress
		FROM 
			Packet p
			JOIN Customer cSender ON p.Sender = cSender.CustomerId
			JOIN Customer cReceiver ON p.Receiver = cReceiver.CustomerId
			JOIN Location lOrigin ON p.Origin = lOrigin.LocationId
			JOIN Location lDestination ON p.Destination = lDestination.LocationId
		WHERE 
			p.Status = ?
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return packets, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(consts.RECEIVED)
	if err != nil {
		return packets, err
	}
	defer rows.Close()

	for rows.Next() {
		var packet dao.PacketDAO
		rows.Scan(
			&packet.PacketId,
			&packet.Weight,
			&packet.Status,
			&packet.Sender.CustomerId,
			&packet.Sender.Name,
			&packet.Sender.Phone,
			&packet.Receiver.CustomerId,
			&packet.Receiver.Name,
			&packet.Receiver.Phone,
			&packet.Origin.LocationId,
			&packet.Origin.LocationName,
			&packet.Origin.Address,
			&packet.Destination.LocationId,
			&packet.Destination.LocationName,
			&packet.Destination.Address,
		)
		packets = append(packets, packet)
	}

	return packets, nil
}

func RetrievePacketDataFromDB(packetID int64) (dao.PacketDAO, error) {

	query := `
		SELECT 
			p.Id AS PacketId,
			p.Weight,
			p.Status,
			cSender.CustomerId AS SenderId,
			cSender.Name AS SenderName,
			cSender.Phone AS SenderPhone,
			cReceiver.CustomerId AS ReceiverId,
			cReceiver.Name AS ReceiverName,
			cReceiver.Phone AS ReceiverPhone,
			lOrigin.LocationId AS OriginId,
			lOrigin.Name AS OriginName,
			lOrigin.Address AS OriginAddress,
			lDestination.LocationId AS DestinationId,
			lDestination.Name AS DestinationName,
			lDestination.Address AS DestinationAddress
		FROM 
			Packet p
			JOIN Customer cSender ON p.Sender = cSender.CustomerId
			JOIN Customer cReceiver ON p.Receiver = cReceiver.CustomerId
			JOIN Location lOrigin ON p.Origin = lOrigin.LocationId
			JOIN Location lDestination ON p.Destination = lDestination.LocationId
		WHERE 
			p.Id = ?
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return dao.PacketDAO{}, err
	}
	defer stmt.Close()

	var packet dao.PacketDAO
	err = stmt.QueryRow(packetID).Scan(
		&packet.PacketId,
		&packet.Weight,
		&packet.Status,
		&packet.Sender.CustomerId,
		&packet.Sender.Name,
		&packet.Sender.Phone,
		&packet.Receiver.CustomerId,
		&packet.Receiver.Name,
		&packet.Receiver.Phone,
		&packet.Origin.LocationId,
		&packet.Origin.LocationName,
		&packet.Origin.Address,
		&packet.Destination.LocationId,
		&packet.Destination.LocationName,
		&packet.Destination.Address,
	)
	if err != nil {
		return dao.PacketDAO{}, err
	}

	return packet, nil
}
