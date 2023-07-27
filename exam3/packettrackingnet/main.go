package main

import (
	"fmt"
	"log"
	"net/http"
	"packettrackingnet/handlers"
	"packettrackingnet/helpers"
	"packettrackingnet/middlewares"
	"packettrackingnet/repository"
	"packettrackingnet/router"
)

func init() {
	helpers.InitErrorLogger()
}

func init() {
	repository.ConnectDatabaseMySQL()
}
func main() {

	mux := http.NewServeMux()

	mux.Handle(router.SHIPMENTS_SENDERS, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS_RECEIVERS, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS_BULK_CREATE, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS_DOWNLOAD, &handlers.ShipmentHandler{})

	mux.Handle(router.LOCATIONS, &handlers.LocationHandler{})

	mux.Handle(router.PACKETS, &handlers.PacketHandler{})
	mux.Handle(router.PACKETS_RECEIVED, &handlers.PacketHandler{})

	mux.Handle(router.TRACKING, &handlers.TrackingHandler{})
	mux.Handle(router.SERVICES, &handlers.ServiceHandler{})
	mux.Handle(router.SERVICES_NAMES, &handlers.ServiceHandler{})
	// Middlewares setup
	var handler http.Handler = mux
	handler = middlewares.Logging(handler)
	handler = middlewares.HandlerAdvice(handler)

	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: handler,
	}

	fmt.Printf("Server is up and running on %s\n", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		helpers.LogError(err)
		return
	}
	log.Println("Shutdown")

}
