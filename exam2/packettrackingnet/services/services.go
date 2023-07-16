package services

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"packettrackingnet/config/consts"
	"packettrackingnet/domain"
	"packettrackingnet/dto"
	"packettrackingnet/helpers"
	"packettrackingnet/repository"
	"strconv"
	"strings"
)

// CreateShipment function to create shipment and save it to the repository.
// Returns the id of newly created shipment.
func CreateShipment(r *http.Request) (*domain.Shipment, error) {
	packet := make(chan *domain.Packet)
	service := make(chan *domain.Service)
	var shipmentRequest dto.ShipmentRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&shipmentRequest); err != nil {
		return nil, errors.New("invalid data")
	}

	err := validator.New().Struct(shipmentRequest)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	go func() {
		_, tmp := repository.FindPacketById(shipmentRequest.PacketID)
		packet <- tmp
	}()

	go func() {
		_, tmp := repository.FindServiceByName(shipmentRequest.ServiceName)
		service <- tmp
	}()

	packetData := <-packet
	serviceData := <-service

	if serviceData == nil || packetData == nil {
		return nil, errors.New("service or packet doesn't exist")
	}
	shippingCost := packetData.Weight * serviceData.PricePerKilogram
	packetData.Status = consts.ON_PROGRESS
	newShipment := domain.NewShipment(packetData, shippingCost, *serviceData, []*domain.Location{packetData.Origin}, false)
	repository.AddShipment(*newShipment)
	return newShipment, nil
}

// GetAllReceivedPackets function to get all packets that has been delivered.
// It checks the status of the Shipment-received flag.
// Returns slice of Packet
func GetAllReceivedPackets() []domain.Packet {
	shipments, _ := repository.GetAllShipment()
	var results = make([]domain.Packet, 0)

	for _, item := range *shipments {
		if item.IsReceived {
			results = append(results, *item.Packet)
		}
	}

	return results
}

// GetAllPacketsByLocationName function to get all the packets that has been going through
// specified location name.
// It returns a slice of PacketDetails struct.
func GetAllPacketsByLocationName(r *http.Request) []domain.PacketDetails {
	shipments, _ := repository.GetAllShipment()
	var results = make([]domain.PacketDetails, 0)
	query := r.URL.Query()
	locationName := query.Get("locationName")
	for _, ship := range *shipments {
		for _, loc := range ship.CheckPoints {
			if strings.EqualFold(loc.LocationName, locationName) {
				results = append(results, *domain.NewPacketDetails(*ship.Packet, ship.IsReceived))
			}
		}
	}
	return results
}

// GetAllCheckpoints function to get all Location saved in repository.
// Returns slice of Location.
func GetAllCheckpoints() []domain.Location {
	results, _ := repository.GetAllLocations()
	return results
}

// UpdateShipmentCheckpoint function to update checkpoint of shipment.
// It only appends the new checkpoint to the end of the checkpoint's slice.
// Returns error if the Shipment or Location is not found.
func UpdateShipmentCheckpoint(r *http.Request) (*domain.Shipment, error) {
	shipmentChan := make(chan *domain.Shipment)
	LocationChan := make(chan *domain.Location)
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
		_, tmp := repository.FindShipmentById(updateShipmentRequest.ShipmentID)
		shipmentChan <- tmp
	}()

	go func() {
		_, tmp := repository.FindLocationById(updateShipmentRequest.LocationID)
		LocationChan <- tmp
	}()

	shipment := <-shipmentChan
	location := <-LocationChan
	if shipment == nil || location == nil {
		return nil, errors.New("data not found")
	}

	shipment.CheckPoints = append(shipment.CheckPoints, location)

	if len(shipment.CheckPoints) > 0 {
		lastPos := shipment.CheckPoints[len(shipment.CheckPoints)-1]
		if strings.EqualFold(lastPos.LocationId, shipment.Packet.Destination.LocationId) {
			shipment.IsReceived = true
			shipment.Packet.Status = consts.RECEIVED
		}
	}
	return shipment, nil
}

// GetAllShipment function to get all shipments saved in repository.
// Returns slice of Shipment struct.
func GetAllShipment() []domain.Shipment {
	results, _ := repository.GetAllShipment()
	return *results
}

// AddSender function to add new sender into repository.
// Returns the newly created Sender pointer.
func AddSender(r *http.Request) (*domain.Sender, error) {
	var sender domain.Sender
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sender); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(sender)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newSender := domain.NewSender(sender.SenderName, sender.Phone)
	repository.AddSender(*newSender)
	return newSender, nil
}

// AddService function to add new service into repository.
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
	repository.AddService(*newService)
	return newService, nil
}

// AddReceiver function to add new receiver into the repository.
// Returns the newly created Receiver pointer.
func AddReceiver(r *http.Request) (*domain.Receiver, error) {
	var receiver domain.Receiver
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receiver); err != nil {
		return nil, errors.New(err.Error())
	}
	defer r.Body.Close()

	err := validator.New().Struct(receiver)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	newReceiver := domain.NewReceiver(receiver.ReceiverName, receiver.Phone)
	repository.AddReceiver(*newReceiver)
	return newReceiver, nil
}

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

	exist, _ := repository.FindLocationByName(location.LocationName)
	if exist {
		return nil, errors.New("location with that name already exist")
	}

	newLocation := domain.NewLocation(location.LocationName, location.Address)
	repository.AddLocation(newLocation)
	return newLocation, nil
}

// AddPacket function to add the new packet into the repository.
// Returns the newly created Packet pointer.
func AddPacket(r *http.Request) (*domain.Packet, error) {
	senderChan := make(chan *domain.Sender)
	receiverChan := make(chan *domain.Receiver)
	originChan := make(chan *domain.Location)
	destinationChan := make(chan *domain.Location)
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
		_, tmp := repository.FindSenderById(packetRequest.SenderID)
		senderChan <- tmp
	}()

	go func() {
		_, tmp := repository.FindReceiverById(packetRequest.ReceiverID)
		receiverChan <- tmp
	}()

	go func() {
		_, tmp := repository.FindLocationByName(packetRequest.OriginName)
		originChan <- tmp
	}()

	go func() {
		_, tmp := repository.FindLocationByName(packetRequest.DestinationName)
		destinationChan <- tmp
	}()

	sender := <-senderChan
	receiver := <-receiverChan
	origin := <-originChan
	destination := <-destinationChan

	if sender == nil || receiver == nil || origin == nil || destination == nil {
		return nil, errors.New("incomplete data")
	}

	newPacket := domain.NewPacket(*sender, *receiver, origin, destination, packetRequest.Weight)
	newPacket.Status = consts.PENDING
	repository.AddPacket(*newPacket)
	return newPacket, nil
}

// GetAllServices function to get all Services from the repository.
// Returns slice of Services.
func GetAllServices() []domain.Service {
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
func GetServiceByName(r *http.Request) *domain.Service {
	query := r.URL.Query()
	name := query.Get("serviceName")
	_, service := repository.FindServiceByName(name)
	return service
}

// GetShipmentById function to query Shipment from repository by Shipment ID.
// Returns boolean to indicate the existence of Shipment and Shipment pointer.
func GetShipmentById(r *http.Request) (bool, *domain.Shipment) {
	query := r.URL.Query()
	id := query.Get("trackingId")
	return repository.FindShipmentById(id)
}

func UpdateLocationAddress(r *http.Request) (error, *domain.Location) {
	var request dto.UpdateLocationAddressRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		return errors.New(err.Error()), nil
	}
	defer r.Body.Close()

	err := validator.New().Struct(request)
	if err != nil {
		return errors.New(err.Error()), nil
	}

	exist, loc := repository.FindLocationByName(request.LocationName)
	if !exist {
		return errors.New("location doesn't exist"), nil
	}

	loc.Address = request.Address
	return nil, loc
}

func PersistData() {
	repository.PersistData()
}

func InitDatastore() {
	repository.InitDatastore()
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
	packet := make(chan *domain.Packet)
	service := make(chan *domain.Service)
	var shipmentRequest dto.ShipmentRequest

	shipmentRequest.PacketID = data[0]
	shipmentRequest.ServiceName = data[1]

	err := validator.New().Struct(shipmentRequest)
	if err != nil {
		return errors.New(err.Error())
	}

	go func() {
		_, tmp := repository.FindPacketById(shipmentRequest.PacketID)
		packet <- tmp
	}()

	go func() {
		_, tmp := repository.FindServiceByName(shipmentRequest.ServiceName)
		service <- tmp
	}()

	packetData := <-packet
	serviceData := <-service

	if serviceData == nil || packetData == nil {
		return errors.New("service or packet doesn't exist")
	}
	shippingCost := packetData.Weight * serviceData.PricePerKilogram
	newShipment := domain.NewShipment(packetData, shippingCost, *serviceData, []*domain.Location{packetData.Origin}, false)
	repository.AddShipment(*newShipment)
	return nil
}

func DownloadAllShipmentData() (*os.File, error) {
	header := []string{"ShipmentID",
		"PacketID", "Weight", "Sender Name", "Receiver Name", "Origin",
		"Destination", "Current Position", "Shipping Cost", "Status",
		"IsReceived",
	}
	shipments, _ := repository.GetAllShipment()
	records := make([][]string, 0)

	for _, item := range *shipments {
		records = append(records, []string{
			item.GetId(), item.Packet.GetId(), strconv.FormatFloat(item.Packet.Weight, 'f', -1, 64),
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

func TruncateData() error {
	err := helpers.TruncateDatabase()
	if err != nil {
		helpers.LogError(err)
		return errors.New(err.Error())
	}
	InitDatastore()
	return nil
}

func GetLocationByName(r *http.Request) (error, *domain.Location) {
	query := r.URL.Query()
	locationName := query.Get("locationName")
	exist, location := repository.FindLocationByName(locationName)

	if !exist {
		return errors.New("location doesn't exist"), nil
	}

	return nil, location
}
