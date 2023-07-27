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
	"strings"
)

func AddSender(sender domain.Customer) (int64, error) {
	result, err := DB.Exec("INSERT INTO Customer(Name, Phone) VALUES(?,?)", sender.Name, sender.Phone)
	if err != nil {
		panic(err)
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func AddReceiver(receiver domain.Customer) (int64, error) {
	result, err := DB.Exec("INSERT INTO Customer(Name, Phone) VALUES(?,?)", receiver.Name, receiver.Phone)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func AddLocation(location domain.Location) (int64, error) {
	result, err := DB.Exec("INSERT INTO Location(Name, Address) VALUES(?,?)", location.LocationName, location.Address)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func AddPacket(packet domain.Packet) (int64, error) {
	result, err := DB.Exec("INSERT INTO Packet(Sender, Receiver, Origin, Destination, Weight, Status) VALUES(?,?,?,?,?,?)", packet.Sender, packet.Receiver, packet.GetOrigin(), packet.GetDestination(), packet.Weight, packet.Status)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func AddService(service domain.Service) (int64, error) {
	result, err := DB.Exec("INSERT INTO Service(ServiceName, PricePerKilogram) VALUES(?,?)", service.ServiceName, service.PricePerKilogram)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

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

// GetServiceByID retrieves a ServiceDAO from the database based on the given ID.
func GetServiceByID(serviceID int64) (dao.ServiceDAO, error) {
	// Prepare the SQL statement.
	stmt, err := DB.Prepare("SELECT ServiceId, ServiceName, PricePerKilogram FROM Service WHERE ServiceId = ?")
	if err != nil {
		return dao.ServiceDAO{}, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the ServiceDAO struct.
	var service dao.ServiceDAO
	err = stmt.QueryRow(serviceID).Scan(&service.ServiceId, &service.ServiceName, &service.PricePerKilogram)
	if err != nil {
		if err == sql.ErrNoRows {
			return dao.ServiceDAO{}, fmt.Errorf("service with ID %d not found", serviceID)
		}
		return dao.ServiceDAO{}, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return service, nil
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
		location, err := GetLocationByID(checkpoint.LocationId)
		if err != nil {
			return nil, fmt.Errorf("failed to get location data: %v", err)
		}
		locations = append(locations, *location)
	}

	return locations, nil
}

// GetLocationByID retrieves a LocationDAO from the database based on the given LocationID.
func GetLocationByID(locationID int64) (*dao.LocationDAO, error) {
	// Prepare the SQL statement to retrieve location data based on the given LocationID.
	stmt, err := DB.Prepare("SELECT LocationId, Name, Address FROM Location WHERE LocationId = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the LocationDAO struct.
	var location dao.LocationDAO
	err = stmt.QueryRow(locationID).Scan(&location.LocationId, &location.LocationName, &location.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location with ID %d not found", locationID)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return &location, nil
}

// GetServiceByName retrieves a ServiceDAO from the database based on the given ID.
func GetServiceByName(serviceName string) (*dao.ServiceDAO, error) {
	// Prepare the SQL statement.
	stmt, err := DB.Prepare("SELECT ServiceId, ServiceName, PricePerKilogram FROM Service WHERE UPPER(serviceName) = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the ServiceDAO struct.
	var service dao.ServiceDAO
	err = stmt.QueryRow(strings.ToUpper(serviceName)).Scan(&service.ServiceId, &service.ServiceName, &service.PricePerKilogram)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service with name %s not found", serviceName)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return &service, nil
}

// FindCustomerByID retrieves a CustomerDAO from the database based on the given CustomerID.
func FindCustomerByID(customerID int64) (*dao.CustomerDAO, error) {
	// Prepare the SQL statement to retrieve customer data based on the given CustomerID.
	stmt, err := DB.Prepare("SELECT CustomerId, Name, Phone FROM Customer WHERE CustomerId = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the CustomerDAO struct.
	var customer dao.CustomerDAO
	err = stmt.QueryRow(customerID).Scan(&customer.CustomerId, &customer.Name, &customer.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer with ID %d not found", customerID)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return &customer, nil
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

// FindLocationByName retrieves a LocationDAO from the database based on the given LocationName.
func FindLocationByName(locationName string) (*dao.LocationDAO, error) {
	// Prepare the SQL statement to retrieve location data based on the given LocationID.
	stmt, err := DB.Prepare("SELECT LocationId, Name, Address FROM Location WHERE upper(Name)= ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the LocationDAO struct.
	var location dao.LocationDAO
	err = stmt.QueryRow(strings.ToUpper(locationName)).Scan(&location.LocationId, &location.LocationName, &location.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location with ID %d not found", locationName)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return &location, nil
}

// FindLocationById retrieves a LocationDAO from the database based on the given LocationName.
func FindLocationById(locationId int64) (*dao.LocationDAO, error) {
	// Prepare the SQL statement to retrieve location data based on the given LocationID.
	stmt, err := DB.Prepare("SELECT LocationId, Name, Address FROM Location WHERE LocationId = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the LocationDAO struct.
	var location dao.LocationDAO
	err = stmt.QueryRow(locationId).Scan(&location.LocationId, &location.LocationName, &location.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location with ID %d not found", locationId)
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return &location, nil
}

// GetAllServices retrieves a ServiceDAO from the database.
func GetAllServices() ([]dao.ServiceDAO, error) {
	var services = make([]dao.ServiceDAO, 0)
	// Prepare the SQL statement.
	stmt, err := DB.Prepare("SELECT ServiceId, ServiceName, PricePerKilogram FROM Service")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the ServiceDAO struct.
	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(err.Error())
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var service dao.ServiceDAO
		rows.Scan(&service.ServiceId, &service.ServiceName, &service.PricePerKilogram)
		services = append(services, service)
	}

	return services, nil

}

// GetAllLocations retrieves a LocationDAO from the database.
func GetAllLocations() ([]dao.LocationDAO, error) {
	var locations = make([]dao.LocationDAO, 0)
	// Prepare the SQL statement.
	stmt, err := DB.Prepare("SELECT LocationId, LocationName, Address FROM Location")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into the LocationDAO struct.
	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(err.Error())
		}
		return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var location dao.LocationDAO
		rows.Scan(&location.LocationId, &location.LocationName, &location.Address)
		locations = append(locations, location)
	}

	return locations, nil

}

func UpdateLocationAddressByLocationName(locationName string) error {
	_, err := DB.Exec("UPDATE Location SET Address = ? WHERE lower(LocationName) = ?", strings.ToLower(locationName))
	if err != nil {
		return err
	}
	return nil
}

func UpdateShipmentStatus(shipmentId int64, status bool) error {
	_, err := DB.Exec("UPDATE Shipment SET IsReceived = ? WHERE ShipmentId = ?", status, shipmentId)
	if err != nil {
		return err
	}
	return nil
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
