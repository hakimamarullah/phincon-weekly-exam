package handlers

import (
	"net/http"
	"packettrackingnet/dto"
	"packettrackingnet/helpers"
	"packettrackingnet/services"
)

func postSender(w http.ResponseWriter, r *http.Request) {
	data, err := services.AddSender(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: data, Code: http.StatusCreated, Count: 1})
}

func postReceiver(w http.ResponseWriter, r *http.Request) {
	data, err := services.AddReceiver(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: data, Code: http.StatusCreated, Count: 1})
}

func postPacket(w http.ResponseWriter, r *http.Request) {
	data, err := services.AddPacket(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: data, Code: http.StatusCreated, Count: 1})
}

func postLocation(w http.ResponseWriter, r *http.Request) {
	data, err := services.AddLocation(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: data, Code: http.StatusCreated, Count: 1})
}

func postService(w http.ResponseWriter, r *http.Request) {
	data, err := services.AddService(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: data, Code: http.StatusCreated, Count: 1})
}

func postShipment(w http.ResponseWriter, r *http.Request) {
	shipment, err := services.CreateShipment(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: shipment, Count: 1, Code: http.StatusCreated})
}

func getServices(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllServices()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func getLocations(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllCheckpoints()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func getAllPacketByLocationName(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllPacketsByLocationName(r)
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func getAllShipments(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllShipment()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func getAllServiceNames(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllServiceNames()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func getServiceByName(w http.ResponseWriter, r *http.Request) {
	result := services.GetServiceByName(r)
	if result == nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: "Not Found", Code: http.StatusNotFound})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: *result, Count: 1})
}

func getShipmentById(w http.ResponseWriter, r *http.Request) {
	_, result := services.GetShipmentById(r)
	if result == nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: "Not Found", Code: http.StatusNotFound})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: *result, Count: 1})
}

func updateShipmentCheckpoint(w http.ResponseWriter, r *http.Request) {
	checkpoint, err := services.UpdateShipmentCheckpoint(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	helpers.ResponseJSON(w, dto.ResponseBody{Data: *checkpoint})
}

func getAllReceivedPackets(w http.ResponseWriter, r *http.Request) {
	results := services.GetAllReceivedPackets()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}

func updateLocationAddressByName(w http.ResponseWriter, r *http.Request) {
	err, updatedLocation := services.UpdateLocationAddress(r)
	if err != nil {
		helpers.ResponseJSON(w, dto.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	helpers.ResponseJSON(w, dto.ResponseBody{Data: *updatedLocation})
}
