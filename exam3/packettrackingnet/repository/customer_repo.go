package repository

import (
	"database/sql"
	"fmt"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
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
