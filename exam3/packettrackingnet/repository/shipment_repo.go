package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"packettrackingnet/config/consts"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
	"packettrackingnet/dto"
)

func AddShipmentCheckpoint(shipmentId, locationId int64) error {
	_, err := DB.Exec("INSERT INTO Checkpoint(ShipmentId, LocationId) VALUES(?,?)", shipmentId, locationId)
	if err != nil {
		return err
	}
	return nil
}
func AddShipment(shipment dto.ShipmentRequest) (int64, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	packet, err := RetrievePacketDataFromDB(shipment.PacketID)
	if err != nil {
		return 0, err
	}
	service, err := GetServiceByName(shipment.ServiceName)
	if err != nil {
		return 0, err
	}

	shippingCost := packet.Weight * service.PricePerKilogram
	result, err := tx.Exec("INSERT INTO Shipment(Packet, Service, ShippingCost) VALUES(?,?,?)", shipment.PacketID, shipment.ServiceId, shippingCost)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	shipmentId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	_, err = tx.Exec("INSERT INTO Checkpoint(ShipmentId, LocationId) VALUES(?,?)", shipmentId, packet.Origin.LocationId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	log.Println(packet.PacketId)
	_, err = tx.Exec("UPDATE Packet SET Status = ? WHERE Id = ?", consts.ON_PROGRESS, packet.PacketId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return shipmentId, nil
}

func GetAllShipments() ([]dao.ShipmentDAO, error) {
	shipments := make([]dao.ShipmentDAO, 0)
	shipStmt, err := DB.Prepare("SELECT ShipmentId, Packet, ShippingCost, Service, IsReceived FROM Shipment")
	if err != nil {
		return shipments, err
	}
	defer shipStmt.Close()

	rows, err := shipStmt.Query()
	if err != nil {
		return shipments, err
	}

	for rows.Next() {
		var shipment dao.ShipmentDAO
		var packetId, serviceId int64
		err := rows.Scan(&shipment.ShipmentId, &packetId, &shipment.ShippingCost, &serviceId, &shipment.IsReceived)
		if err != nil {
			return shipments, err
		}
		packet, _ := RetrievePacketDataFromDB(packetId)
		service, _ := GetServiceByID(serviceId)
		checkpoints, _ := GetAllCheckpointsByShipmentID(shipment.ShipmentId)
		shipment.Packet = packet
		shipment.Service = service
		shipment.CheckPoints = checkpoints
		shipments = append(shipments, shipment)

	}
	return shipments, nil

}

// GetAllCheckpointsByShipmentID retrieves all checkpoints of a shipment from the database based on the given shipment ID.
func GetAllCheckpointsByShipmentID(shipmentID int64) ([]dao.LocationDAO, error) {
	// Prepare the SQL statement to retrieve checkpoint data based on the given shipment ID.
	stmt, err := DB.Prepare("SELECT ShipmentId, LocationId FROM Checkpoint WHERE ShipmentId = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the results into a slice of Checkpoint structs.
	rows, err := stmt.Query(shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}
	defer rows.Close()
	type Checkpoint struct {
		ShipmentId int64
		LocationId int64
	}
	var checkpoints []Checkpoint
	for rows.Next() {
		var checkpoint Checkpoint
		err = rows.Scan(&checkpoint.ShipmentId, &checkpoint.LocationId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan checkpoint data: %v", err)
		}
		checkpoints = append(checkpoints, checkpoint)
	}

	// Now fetch the location data for each checkpoint.
	var locations []dao.LocationDAO
	for _, checkpoint := range checkpoints {
		location, err := FindLocationById(checkpoint.LocationId)
		if err != nil {
			return nil, fmt.Errorf("failed to get location data: %v", err)
		}
		locations = append(locations, *location)
	}

	return locations, nil
}

// FindShipmentById retrieves a ShipmentDAO from the database based on the given ShipmentID.
func FindShipmentById(shipmentId int64) (*dao.ShipmentDAO, error) {
	// Prepare the SQL statement to retrieve customer data based on the given CustomerID.
	stmt, err := DB.Prepare("SELECT ShipmentId, Packet, ShippingCost, Service, IsReceived FROM Shipment where ShipmentId = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the ShipmentDAO struct.
	var shipment domain.Shipment
	err = stmt.QueryRow(shipmentId).Scan(&shipment.ShipmentId, &shipment.Packet, &shipment.ShippingCost, &shipment.Service, &shipment.IsReceived)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shipment with ID %d not found", shipmentId)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	packet, err := RetrievePacketDataFromDB(shipment.Packet)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	service, err := GetServiceByID(shipment.Service)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	checkpoints, err := GetAllCheckpointsByShipmentID(shipmentId)
	if err != nil {
		return nil, err
	}
	result := dao.NewShipmentDAO(shipmentId, packet, shipment.ShippingCost, service, checkpoints, shipment.IsReceived)
	return result, nil
}

func UpdateShipmentStatus(shipmentId int64, status bool) error {
	_, err := DB.Exec("UPDATE Shipment SET IsReceived = ? WHERE ShipmentId = ?", status, shipmentId)
	if err != nil {
		return err
	}
	return nil
}
