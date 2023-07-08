package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"packettracking/domain"
	"packettracking/services"
	"packettracking/utils"
	"strings"
)

func main() {
	fmt.Println("WELCOME TO GO TRACKING")

	services.CreateService("ONS", 2000.0)
	services.CreateService("Sameday", 1500.0)
	services.CreateService("Cargo", 1000.0)

	var serviceNames = services.GetAllServiceNames()

menuLoop:
	for {
		prompt := promptui.Select{Label: "Menu", Items: []string{"Shipment",
			"List Received Packets",
			"List Packet Based On Checkpoint",
			"List Checkpoint",
			"Update Shipment Checkpoint",
			"Add Location",
			"List Shipments",
			"Tracking",
			"List Service",
			"Exit"}}
		cmd, _, _ := prompt.Run()

		switch cmd {
		case 0:
			createShipment(serviceNames)
		case 1:
			displayAllReceivedPackets()
		case 2:
			displayAllPacketDetailsByLocationName()
		case 3:
			displayCheckpoints()
		case 4:
			updateShipmentCheckpoint()
		case 5:
			addNewCheckPoint()
		case 6:
			displayAllShipments()
		case 7:
			tracking()
		case 8:
			displayAllService()
		case 9:
			stop()
			break menuLoop
		}

	}

}

// stop procedure to display a goodbye message
func stop() {
	fmt.Println("Thank You For Using Go Tracking!")
}

// updateShipmentCheckpoint procedure to update Shipment Checkpoint
func updateShipmentCheckpoint() {
	shipmentId := utils.ScanStringNonBlank("Shipment ID")
	locationId := utils.ScanStringNonBlank("Location ID")
	_, err := services.UpdateShipmentCheckpoint(shipmentId, locationId)
	if err != nil {
		fmt.Println(fmt.Errorf("%s\n", err.Error()))
		return
	}
	fmt.Println("Checkpoint Updated!!")
}

// displayCheckpoints procedure to display all available locations in the repository
func displayCheckpoints() {
	checkPoints := services.GetAllCheckpoints()
	t := utils.StandardTable("Checkpoints")

	t.AppendHeader(table.Row{"ID", "Location Name", "Address"})

	for _, item := range checkPoints {
		t.AppendRow(table.Row{item.Id, item.LocationName, item.Address})
	}
	fmt.Println(t.Render())
}

// displayAllReceivedPackets procedure to display all Packets that has been delivered
func displayAllReceivedPackets() {
	t := utils.StandardTable("Delivered Packets")

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Weight (kg)", "Address"})

	for _, item := range services.GetAllReceivedPackets() {
		t.AppendRow(table.Row{item.Id, item.Sender.SenderName, item.Receiver.ReceiverName, item.Weight, item.Destination.Address})
	}

	fmt.Println(t.Render())

}

// createShipment procedure to orchestrating the creation of shipment
func createShipment(serviceNames []string) {
	name := utils.ScanStringNonBlank("Sender Name")
	phone := utils.ScanStringNonBlank("Phone")
	sender := services.AddSender(name, phone)
	fmt.Printf("Sender ID: %s\n", sender.Id)

	name = utils.ScanStringNonBlank("Receiver Name")
	phone = utils.ScanStringNonBlank("Phone")
	receiver := services.AddReceiver(name, phone)

	prompt := promptui.Select{Label: "Service Type", Items: serviceNames}
	_, serviceType, _ := prompt.Run()
	service := services.GetServiceByName(serviceType)

	origin := utils.ScanStringNonBlank("Origin Name")
	address := utils.ScanStringNonBlank("Address")
	_, originLocation := services.AddLocation(origin, address)

	dest := utils.ScanStringNonBlank("Destination Name")
	address = utils.ScanStringNonBlank("Address")
	_, destLocation := services.AddLocation(dest, address)

	weight := utils.ScanInt("Packet Weight")
	packet := services.AddPacket(*sender, *receiver, *originLocation, *destLocation, float64(weight))

	confirm := utils.ConfirmInput("Are you sure to create the shipment?")

	if confirm {
		shipmentId := services.CreateShipment(*packet, *service)
		fmt.Printf("SHIPMENT CREATED!\nSHIPMENT ID: %s\n", strings.ToUpper(shipmentId))
		return
	}
	fmt.Println("CANCELLED!")
}

// addNewCheckPoint procedure to add new location into repository
func addNewCheckPoint() {
	name := utils.ScanStringNonBlank("Location Name")
	address := utils.ScanStringNonBlank("Address")
	added, _ := services.AddLocation(name, address)

	if added {
		fmt.Println("Location added successfully!")
	} else {
		fmt.Println("Location already exist! Not creating!")
	}
}

// displayAllPacketDetailsByLocationName procedure to display all the packet's details,
// which has been going through the checkpoint specified by location name.
func displayAllPacketDetailsByLocationName() {
	locationName := utils.ScanStringNonBlank("Location Name")
	packetDetails := services.GetAllPacketsByLocationName(locationName)

	t := utils.StandardTable(locationName)

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Weight", "Origin", "Destination", "Status"})

	for _, item := range packetDetails {
		packet := item.Packet
		t.AppendRow(table.Row{
			packet.Id,
			packet.Sender.SenderName,
			packet.Receiver.ReceiverName,
			packet.Weight,
			packet.Origin.LocationName,
			packet.Destination.LocationName,
			utils.GetStatus(item.IsReceived),
		})
	}

	fmt.Println(t.Render())
}

// displayAllShipments procedure to display all the shipment in repository
func displayAllShipments() {
	shipments := services.GetAllShipment()

	t := utils.StandardTable("SHIPMENTS")

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Packet ID", "Service", "Total Price", "Destination ID", "Current Pos", "Status"})

	for _, item := range shipments {
		t.AppendRow(table.Row{
			item.Id,
			item.Packet.Sender.SenderName,
			item.Packet.Receiver.ReceiverName,
			item.Packet.Id,
			item.Service.ServiceName,
			item.ShippingCost,
			item.Packet.Destination.Id,
			currentCheckpointId(item.CheckPoints),
			utils.GetStatus(item.IsReceived),
		})
	}

	fmt.Println(t.Render())

}

// tracking procedure to track the packet by using Shipment ID
// this displays the details in table format.
func tracking() {
	shipmentId := utils.ScanStringNonBlank("SHIPMENT ID")
	exist, data := services.GetShipmentById(shipmentId)
	if !exist {
		fmt.Println(fmt.Errorf("%s", "Shipment Not Found"))
		return
	}

	t := utils.StandardTable("TRACKING")

	t.AppendHeader(table.Row{"History"})
	t.AppendHeader(table.Row{fmt.Sprintf("SHIPMENT ID: %s", strings.ToUpper(data.Id))})
	t.AppendHeader(table.Row{fmt.Sprintf("Sender: %s", data.Packet.Sender.SenderName)})
	t.AppendHeader(table.Row{fmt.Sprintf("Receiver: %s", data.Packet.Receiver.ReceiverName)})
	t.AppendHeader(table.Row{fmt.Sprintf("Origin: %s", data.Packet.Origin.Id)})
	t.AppendHeader(table.Row{fmt.Sprintf("Destination: %s", data.Packet.Destination.Id)})

	checkPoints := data.CheckPoints
	for _, item := range checkPoints {
		t.AppendRow(table.Row{
			fmt.Sprintf("ID: %s\nLocation Name: %s\nAddress: %s\n", item.Id, item.LocationName, strings.ToTitle(item.Address)),
		})
	}

	fmt.Println(t.Render())
}

// currentCheckpointId function to get the id of the last checkpoint of the shipment.
// It returns the string Location ID of the last checkpoint
func currentCheckpointId(checkPoints []domain.Location) string {
	if len(checkPoints) == 0 {
		return "NONE"
	}
	lastIndex := len(checkPoints) - 1
	return checkPoints[lastIndex].Id
}

// displayAllService procedure to display all available shipment service types.
func displayAllService() {
	serviceList := services.GetAllServices()

	t := utils.StandardTable("Service Type")

	t.AppendHeader(table.Row{"Service Name", "Price per Kilogram"})

	for _, item := range serviceList {
		t.AppendRow(table.Row{
			item.ServiceName,
			item.PricePerKilogram,
		})
	}

	fmt.Println(t.Render())
}
