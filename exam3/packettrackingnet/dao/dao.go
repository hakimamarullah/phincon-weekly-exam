package dao

type CustomerDAO struct {
	CustomerId int64  `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

func NewCustomerDAO(id int64, name string, phone string) *CustomerDAO {
	return &CustomerDAO{Name: name, Phone: phone}
}

type PacketDAO struct {
	PacketId    int64       `json:"packetId" validate:"required"`
	Sender      CustomerDAO `json:"sender" validate:"required"`
	Receiver    CustomerDAO `json:"receiver" validate:"required"`
	Origin      LocationDAO `json:"origin" validate:"required"`
	Destination LocationDAO `json:"destination" validate:"required"`
	Weight      float64     `json:"weight"`
	Status      string      `json:"status"`
}

func NewPacketDAO(packetId int64, sender CustomerDAO, receiver CustomerDAO, origin LocationDAO, destination LocationDAO, weight float64, status string) *PacketDAO {
	return &PacketDAO{PacketId: packetId, Sender: sender, Receiver: receiver, Origin: origin, Destination: destination, Weight: weight, Status: status}
}

type ShipmentDAO struct {
	ShipmentId   int64         `json:"shipmentId"`
	Packet       PacketDAO     `validate:"required"`
	ShippingCost float64       `json:"shippingCost" validate:"required"`
	Service      ServiceDAO    `validate:"required"`
	CheckPoints  []LocationDAO `json:"checkPoints"`
	IsReceived   bool          `json:"isReceived"`
}

func NewShipmentDAO(shipmentId int64, packet PacketDAO, shippingCost float64, service ServiceDAO, checkPoints []LocationDAO, isReceived bool) *ShipmentDAO {
	return &ShipmentDAO{ShipmentId: shipmentId, Packet: packet, ShippingCost: shippingCost, Service: service, CheckPoints: checkPoints, IsReceived: isReceived}
}

type LocationDAO struct {
	LocationId   int64  `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address" validate:"required"`
}

func NewLocationDAO(locationId int64, locationName string, address string) *LocationDAO {
	return &LocationDAO{LocationId: locationId, LocationName: locationName, Address: address}
}

type ServiceDAO struct {
	ServiceId        int64   `json:"serviceId"`
	ServiceName      string  `json:"serviceName" validate:"required"`
	PricePerKilogram float64 `json:"pricePerKilogram" validate:"required"`
}

func NewServiceDAO(serviceId int64, serviceName string, pricePerKilogram float64) *ServiceDAO {
	return &ServiceDAO{ServiceId: serviceId, ServiceName: serviceName, PricePerKilogram: pricePerKilogram}
}

type PacketDetails struct {
	Packet     PacketDAO
	IsReceived bool `json:"isReceived"`
}

func NewPacketDetails(packet PacketDAO, isReceived bool) *PacketDetails {
	return &PacketDetails{Packet: packet, IsReceived: isReceived}
}

func (sh *ShipmentDAO) GetId() int {
	return int(sh.ShipmentId)
}

func (sh *ShipmentDAO) GetCurrentPosition() string {
	return sh.CheckPoints[len(sh.CheckPoints)-1].LocationName
}

func (pk *PacketDAO) GetId() int {
	return int(pk.PacketId)
}

func (pk *PacketDAO) GetSender() string {
	return pk.Sender.Name
}

func (pk *PacketDAO) GetReceiver() string {
	return pk.Receiver.Name
}

func (pk *PacketDAO) GetOrigin() string {
	return pk.Origin.LocationName
}

func (pk *PacketDAO) GetDestination() string {
	return pk.Destination.LocationName
}
