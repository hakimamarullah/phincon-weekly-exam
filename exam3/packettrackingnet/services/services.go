package services

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"packettrackingnet/config/consts"
	"packettrackingnet/dao"
	"packettrackingnet/domain"
	"packettrackingnet/dto"
	"packettrackingnet/helpers"
	"packettrackingnet/mapper"
	"packettrackingnet/repository"
	"strconv"
	"strings"
	"sync"
)

// CreateShipment function to create shipment and save it to the repository.
// Returns the id of newly created shipment.
func CreateShipment(r *http.Request) (int64, error) {
	var shipmentRequest dto.ShipmentRequest

	if err := json.NewDecoder(r.Body).Decode(&shipmentRequest); err != nil {
		return 0, errors.New(err.Error())
	}
	defer r.Body.Close()
	service, err := repository.GetServiceByName(shipmentRequest.ServiceName)
	if err != nil {
		return 0, err
	}
	shipmentRequest.ServiceId = service.ServiceId
	serviceId, err := repository.AddShipment(shipmentRequest)
	if err != nil {
		return 0, err
	}
	return serviceId, nil
}

// GetAllReceivedPackets function to get all packets that has been delivered.
// It checks the status of the Shipment-received flag.
// Returns slice of Packet
func GetAllReceivedPackets() []dao.PacketDAO {
	packets, err := repository.GetAllReceivedPackets()
	if err != nil {
		return nil
	}
	return packets
}

// GetAllPacketsByLocationName function to get all the packets that has been going through
// specified location name.
// It returns a slice of PacketDetails struct.
func GetAllPacketsByLocationName(r *http.Request) []dao.PacketDetails {
	shipments, _ := repository.GetAllShipments()
	var results = make([]dao.PacketDetails, 0)
	query := r.URL.Query()
	locationName := query.Get("locationName")
	for _, ship := range shipments {
		for _, loc := range ship.CheckPoints {
			if strings.EqualFold(loc.LocationName, locationName) {
				results = append(results, *dao.NewPacketDetails(ship.Packet, ship.IsReceived))
			}
		}
	}
	return results
}

// GetAllCheckpoints function to get all Location saved in repository.
// Returns slice of Location.
func GetAllCheckpoints() []dao.LocationDAO {
	results, _ := repository.GetAllLocations()
	return results
}

// UpdateShipmentCheckpoint function to update checkpoint of shipment.
// It only appends the new checkpoint to the end of the checkpoint's slice.
// Returns error if the Shipment or Location is not found.
func UpdateShipmentCheckpoint(r *http.Request) (*dto.ShipmentResponse, error) {
	shipmentChan := make(chan *dao.ShipmentDAO)
	LocationChan := make(chan *dao.LocationDAO)
	var updateShipmentRequest dto.UpdateShipmentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateShipmentRequest); err != nil {
		return nil, errors.New("invalid data")
	}

	err := validator.New().Struct(updateShipmentRequest)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	go func() {
		tmp, _ := repository.FindShipmentById(updateShipmentRequest.ShipmentID)
		shipmentChan <- tmp
	}()

	go func() {
		tmp, _ := repository.FindLocationById(updateShipmentRequest.LocationID)
		LocationChan <- tmp
	}()

	shipment := <-shipmentChan
	location := <-LocationChan
	if shipment == nil || location == nil {
		return nil, errors.New("data not found")
	}

	for _, item := range shipment.CheckPoints {
		if item.LocationId == location.LocationId {
			return mapper.ShipmentDaoToShipmentResponse(*shipment), errors.New("duplicate checkpoint")
		}
	}

	shipment.CheckPoints = append(shipment.CheckPoints, *location)
	err = repository.AddShipmentCheckpoint(shipment.ShipmentId, location.LocationId)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	if len(shipment.CheckPoints) > 0 {
		lastPos := shipment.CheckPoints[len(shipment.CheckPoints)-1]
		if lastPos.LocationId == shipment.Packet.Destination.LocationId {
			go func() {
				repository.UpdateShipmentStatus(shipment.ShipmentId, true)
				wg.Done()
			}()

			go func() {
				repository.UpdatePacketStatus(shipment.Packet.PacketId, consts.RECEIVED)
				wg.Done()
			}()
		}
	}
	updatedShipment, _ := repository.FindShipmentById(shipment.ShipmentId)

	return mapper.ShipmentDaoToShipmentResponse(*updatedShipment), nil
}

// GetAllShipment function to get all shipments saved in repository.
// Returns slice of Shipment struct.
func GetAllShipment() []dto.ShipmentResponse {
	shipments, _ := repository.GetAllShipments()
	response := make([]dto.ShipmentResponse, 0)
	for _, item := range shipments {
		response = append(response, *mapper.ShipmentDaoToShipmentResponse(item))
	}
	return response
}

// AddSender function to add new sender into repository.
// Returns the newly created Sender pointer.
func AddSender(r *http.Request) (*domain.Customer, error) {
	var sender domain.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sender); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(sender)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newSender := domain.NewCustomer(sender.Name, sender.Phone)
	id, err := repository.AddSender(*newSender)
	if err != nil {
		return nil, err
	}
	newSender.Id = id
	return newSender, nil
}

// AddReceiver function to add new sender into repository.
// Returns the newly created Sender pointer.
func AddReceiver(r *http.Request) (*domain.Customer, error) {
	var receiver domain.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receiver); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(receiver)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newReceiver := domain.NewCustomer(receiver.Name, receiver.Phone)
	id, err := repository.AddReceiver(*newReceiver)
	if err != nil {
		return nil, err
	}
	newReceiver.Id = id
	return newReceiver, nil
}

//

// AddService function to add new service into the repository.
// Returns the newly created Service pointer.
func AddService(r *http.Request) (*domain.Service, error) {
	var service domain.Service
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&service); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(service)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newService := domain.NewService(service.ServiceName, service.PricePerKilogram)
	id, err := repository.AddService(*newService)
	if err != nil {
		return nil, err
	}
	newService.ServiceId = id
	return newService, nil
}

//

// AddLocation function to add new Location into the repository.
// Returns the newly created Location pointer.
// It only puts the Location into the repository if the new Location name doesn't exist in the repository.
func AddLocation(r *http.Request) (*domain.Location, error) {
	var location domain.Location
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&location); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(location)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newLocation := domain.NewLocation(location.LocationName, location.Address)
	id, err := repository.AddLocation(*newLocation)
	if err != nil {
		return nil, err
	}
	newLocation.LocationId = id
	return newLocation, nil
}

// AddPacket function to add the new packet into the repository.
// Returns the newly created Packet pointer.
func AddPacket(r *http.Request) (*dao.PacketDAO, error) {
	senderChan := make(chan *dao.CustomerDAO)
	receiverChan := make(chan *dao.CustomerDAO)
	originChan := make(chan *dao.LocationDAO)
	destinationChan := make(chan *dao.LocationDAO)
	var packetRequest dto.PacketRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&packetRequest); err != nil {
		return nil, errors.New("invalid data")
	}

	err := validator.New().Struct(packetRequest)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	go func() {
		tmp, _ := repository.FindCustomerByID(packetRequest.SenderID)
		senderChan <- tmp
	}()

	go func() {
		tmp, _ := repository.FindCustomerByID(packetRequest.ReceiverID)
		receiverChan <- tmp
	}()

	go func() {
		tmp, _ := repository.FindLocationByName(packetRequest.OriginName)
		originChan <- tmp
	}()

	go func() {
		tmp, _ := repository.FindLocationByName(packetRequest.DestinationName)
		destinationChan <- tmp
	}()

	sender := <-senderChan
	receiver := <-receiverChan
	origin := <-originChan
	destination := <-destinationChan

	if sender == nil || receiver == nil || origin == nil || destination == nil {
		return nil, errors.New("incomplete data")
	}

	newPacket := domain.NewPacket(sender.CustomerId, receiver.CustomerId, origin.LocationId, destination.LocationId, packetRequest.Weight)
	newPacket.Status = consts.PENDING
	id, err := repository.AddPacket(*newPacket)
	if err != nil {
		return nil, err
	}
	newPacket.PacketId = id
	packet, _ := repository.RetrievePacketDataFromDB(id)
	return &packet, nil
}

// GetAllServices function to get all Services from the repository.
// Returns slice of Services.
func GetAllServices() []dao.ServiceDAO {
	results, _ := repository.GetAllServices()
	return results
}

// GetAllServiceNames function to get all service names from repository,
// Returns slice of service's name
func GetAllServiceNames() []string {
	var results []string
	services, _ := repository.GetAllServices()
	for _, item := range services {
		results = append(results, item.ServiceName)
	}
	return results
}

// GetServiceByName function to query service from repository by service name.
// This method is case-insensitive.
// Returns a single Service pointer.
func GetServiceByName(r *http.Request) *dao.ServiceDAO {
	query := r.URL.Query()
	name := query.Get("serviceName")
	service, _ := repository.GetServiceByName(name)
	return service
}

// GetShipmentById function to query Shipment from repository by Shipment ID.
// Returns boolean to indicate the existence of Shipment and Shipment pointer.
func GetShipmentById(r *http.Request) (bool, *dto.ShipmentResponse) {
	query := r.URL.Query()
	id := query.Get("trackingId")
	intId, err2 := strconv.Atoi(id)
	if err2 != nil {
		return false, nil
	}
	shipment, err := repository.FindShipmentById(int64(intId))
	if err != nil {
		return false, nil
	}
	return true, mapper.ShipmentDaoToShipmentResponse(*shipment)
}

func UpdateLocationAddress(r *http.Request) error {
	var request dto.UpdateLocationAddressRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		return errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(request)
	if err != nil {
		return errors.New(err.Error())
	}

	err = repository.UpdateLocationAddressByLocationName(request.LocationName)
	if err != nil {
		return errors.New("location doesn't exist")
	}

	return nil
}

func UploadShipmentCSV(r *http.Request) error {
	file, _, err := r.FormFile("file")
	if err != nil {
		return errors.New(err.Error())
	}
	defer file.Close()

	records, err := helpers.ReadUploadedCSV(file, false)
	errorChan := make(chan error)
	for _, record := range records {
		record := record
		go func() {
			err := createShipmentFromCSV(record)
			errorChan <- err
		}()
	}
	err = <-errorChan
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func createShipmentFromCSV(data []string) error {
	var shipmentRequest dto.ShipmentRequest

	intPacketId, _ := strconv.Atoi(data[0])
	shipmentRequest.PacketID = int64(intPacketId)
	shipmentRequest.ServiceName = data[1]

	service, err := repository.GetServiceByName(data[1])
	if err != nil {
		return err
	}
	shipmentRequest.ServiceId = service.ServiceId

	err = validator.New().Struct(shipmentRequest)
	if err != nil {
		return errors.New(err.Error())
	}

	_, err = repository.AddShipment(shipmentRequest)
	if err != nil {
		return err
	}
	return nil
}

func DownloadAllShipmentData() (*os.File, error) {
	header := []string{"ShipmentID",
		"PacketID", "Weight", "Sender Name", "Receiver Name", "Origin",
		"Destination", "Current Position", "Shipping Cost", "Status",
		"IsReceived",
	}
	shipments, _ := repository.GetAllShipments()
	records := make([][]string, 0)

	for _, item := range shipments {
		records = append(records, []string{
			strconv.Itoa(item.GetId()), strconv.Itoa(item.Packet.GetId()), strconv.FormatFloat(item.Packet.Weight, 'f', -1, 64),
			item.Packet.GetSender(), item.Packet.GetReceiver(),
			item.Packet.GetOrigin(), item.Packet.GetDestination(), item.GetCurrentPosition(),
			strconv.FormatFloat(item.ShippingCost, 'f', -1, 64),
			item.Packet.Status, strconv.FormatBool(item.IsReceived),
		})
	}

	file, err := helpers.WriteCSV(header, records)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return file, nil

}

func GetLocationByName(r *http.Request) (error, *dao.LocationDAO) {
	query := r.URL.Query()
	locationName := query.Get("locationName")
	location, err := repository.FindLocationByName(locationName)

	if err != nil {
		return errors.New("location doesn't exist"), nil
	}

	return nil, location
}
