package main

import (
	"fmt"
	"log"
	"net/http"
	"packettrackingnet/handlers"
	"packettrackingnet/helpers"
	"packettrackingnet/middlewares"
	"packettrackingnet/router"
)

func init() {
	err := helpers.TruncateDatabase()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
func main() {
	mux := http.NewServeMux()
	mux.Handle(router.SHIPMENTS_SENDERS, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS_RECEIVERS, &handlers.ShipmentHandler{})
	mux.Handle(router.SHIPMENTS, &handlers.ShipmentHandler{})
	mux.Handle(router.LOCATIONS, &handlers.LocationHandler{})
	mux.Handle(router.PACKETS, &handlers.PacketHandler{})
	mux.Handle(router.TRACKING, &handlers.TrackingHandler{})
	mux.Handle(router.SERVICES, &handlers.ServiceHandler{})

	// Middlewares setup
	var handler http.Handler = mux
	handler = middlewares.Logging(handler)
	//handler = middlewares.Authenticate(handler)
	handler = middlewares.HandlerAdvice(handler)

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: handler,
	}

	fmt.Printf("Server is up and running on %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

}
