package handlers

import (
	"net/http"
	"packettrackingnet/dto"
	"packettrackingnet/helpers"
	"packettrackingnet/repository"
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
	results, _ := repository.GetAllServices()
	helpers.ResponseJSON(w, dto.ResponseBody{Data: results, Count: len(results)})
}
