package repository

import (
	"database/sql"
	"fmt"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
	"strings"
)

func AddLocation(location domain.Location) (int64, error) {
	result, err := DB.Exec("INSERT INTO Location(Name, Address) VALUES(?,?)", location.LocationName, location.Address)
	if err != nil {
		return 0, err
	}
	lastId, _ := result.LastInsertId()
	return lastId, nil
}

func UpdateLocationAddressByLocationName(locationName string) error {
	_, err := DB.Exec("UPDATE Location SET Address = ? WHERE lower(LocationName) = ?", strings.ToLower(locationName))
	if err != nil {
		return err
	}
	return nil
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
