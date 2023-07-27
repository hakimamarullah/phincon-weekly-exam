package repository

import (
	"database/sql"
	"fmt"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
	"strings"
)

func AddService(service domain.Service) (int64, error) {
	result, err := DB.Exec("INSERT INTO Service(ServiceName, PricePerKilogram) VALUES(?,?)", service.ServiceName, service.PricePerKilogram)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
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
