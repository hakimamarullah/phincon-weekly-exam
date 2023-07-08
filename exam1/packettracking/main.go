package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

		if cmd == 0 {
			createShipment(serviceNames)
		}

		if cmd == 1 {
			displayAllReceivedPackets()
		}

		if cmd == 2 {
			displayALlPacketsDetailsByLocationName()
		}

		if cmd == 3 {
			displayCheckpoints()
		}

		if cmd == 4 {
			updateShipmentCheckpoint()
		}
		if cmd == 5 {
			addNewCheckPoint()
		}
		if cmd == 6 {
			displayAllShipments()
		}
		if cmd == 7 {
			tracking()
		}
		if cmd == 8 {
			displayAllService()
		}
		if cmd == 9 {
			break
		}
	}

}

func updateShipmentCheckpoint() {
	shipmentId := utils.ScanStringNonBlank("Shipment ID")
	locationId := utils.ScanStringNonBlank("Location ID")
	data, err := services.UpdateShipmentCheckpoint(shipmentId, locationId)
	fmt.Println(data)
	if err != nil {
		fmt.Println(fmt.Errorf("%s\n", err.Error()))
		return
	}
	fmt.Println("Checkpoint Updated!!")
}

func displayCheckpoints() {
	checkPoints := services.GetAllCheckpoints()
	t := table.NewWriter()
	t.SetTitle("Checkpoints")
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Location Name", "Address"})

	for _, item := range checkPoints {
		t.AppendRow(table.Row{item.Id, item.LocationName, item.Address})
	}
	fmt.Println(t.Render())
}

func displayAllReceivedPackets() {
	t := table.NewWriter()
	t.SetTitle("Delivered Packets")
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Weight", "Address"})

	for _, item := range services.GetAllReceivedPackets() {
		t.AppendRow(table.Row{item.Id, item.Sender.SenderName, item.Receiver.ReceiverName, item.Weight, item.Destination.Address})
	}

	fmt.Println(t.Render())

}

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

	name = utils.ScanStringNonBlank("Location Name")

	address := utils.ScanStringNonBlank("Address")
	destination := services.AddLocation(name, address)
	fmt.Println(destination)

	weight := utils.ScanInt("Packet Weight")
	packet := services.AddPacket(*sender, *receiver, *destination, float64(weight))

	confirm := utils.ConfirmInput("Are you sure to create the shipment?")

	if confirm {
		shipmentId := services.CreateShipment(*packet, *service)
		fmt.Printf("SHIPMENT CREATED!\nSHIPMENT ID: %s\n", strings.ToUpper(shipmentId))
		return
	}
	fmt.Println("CANCELLED!")
}

func addNewCheckPoint() {
	name := utils.ScanStringNonBlank("Location Name")
	address := utils.ScanStringNonBlank("Address")
	services.AddLocation(name, address)
	fmt.Println("Location added successfully!")
}

func displayALlPacketsDetailsByLocationName() {
	locationName := utils.ScanStringNonBlank("Location Name")
	packetDetails := services.GetAllPacketsByLocationName(locationName)

	t := table.NewWriter()
	t.SetTitle(strings.ToTitle(locationName))
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Weight", "Address", "Status"})

	for _, item := range packetDetails {
		packet := item.Packet
		t.AppendRow(table.Row{
			packet.Id,
			packet.Sender.SenderName,
			packet.Receiver.ReceiverName,
			packet.Weight,
			packet.Destination.Address,
			utils.GetStatus(item.IsReceived),
		})
	}

	fmt.Println(t.Render())
}

func displayAllShipments() {
	shipments := services.GetAllShipment()

	t := table.NewWriter()
	t.SetTitle(strings.ToTitle("SHIPMENTS"))
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Packet ID", "Service", "Total Price", "Destination ID", "Current Pos", "Status"})

	for _, item := range shipments {
		fmt.Println(item.CheckPoints)
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

func tracking() {
	shipmentId := utils.ScanStringNonBlank("SHIPMENT ID")
	exist, data := services.GetShipmentById(shipmentId)
	if !exist {
		fmt.Println(fmt.Errorf("%s", "Shipment Not Found"))
		return
	}

	t := table.NewWriter()
	t.SetTitle(strings.ToTitle("TRACKING"))
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Sender", "Receiver", "Packet ID", "Service", "Total Price", "Destination ID", "Current Pos", "Status"})
	t.AppendRow(table.Row{
		data.Id,
		data.Packet.Sender.SenderName,
		data.Packet.Receiver.ReceiverName,
		data.Packet.Id,
		data.Service.ServiceName,
		data.ShippingCost,
		data.Packet.Destination.Id,
		currentCheckpointId(data.CheckPoints),
		utils.GetStatus(data.IsReceived),
	})

	fmt.Println(t.Render())
}

func currentCheckpointId(checkPoints []domain.Location) string {
	if len(checkPoints) == 0 {
		return "NONE"
	}
	lastIndex := len(checkPoints) - 1
	return checkPoints[lastIndex].Id
}

func displayAllService() {
	serviceList := services.GetAllServices()

	t := table.NewWriter()
	t.SetTitle(strings.ToTitle("SERVICE TYPE"))
	t.SetAutoIndex(
		true)
	t.Style().Format.Header = text.FormatTitle
	t.Style().Format.Footer = text.FormatTitle
	t.SetStyle(table.StyleBold)

	t.AppendHeader(table.Row{"ID", "Service Name", "Price per Kilogram"})

	for _, item := range serviceList {
		t.AppendRow(table.Row{
			item.Id,
			item.ServiceName,
			item.PricePerKilogram,
		})
	}

	fmt.Println(t.Render())
}
