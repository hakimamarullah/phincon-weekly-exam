package handlers

import (
	"net/http"
	"packettrackingnet/router"
	"strings"
)

type ShipmentHandler struct{}
type LocationHandler struct{}
type PacketHandler struct{}
type TrackingHandler struct{}
type ServiceHandler struct{}

func (sh *ShipmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.SHIPMENTS_SENDERS) && r.Method == http.MethodPost:
		postSender(w, r)
	case strings.EqualFold(path, router.SHIPMENTS_RECEIVERS) && r.Method == http.MethodPost:
		postReceiver(w, r)
	case strings.EqualFold(path, router.SHIPMENTS) && r.Method == http.MethodPost:
		postShipment(w, r)
	}
}

func (lh *LocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.LOCATIONS) && r.Method == http.MethodPost:
		postLocation(w, r)
	}
}

func (ph *PacketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.PACKETS) && r.Method == http.MethodPost:
		postPacket(w, r)
	}
}

func (th *TrackingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (svh *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.SERVICES) && r.Method == http.MethodPost:
		postService(w, r)
	case strings.EqualFold(path, router.SERVICES) && r.Method == http.MethodGet:
		getServices(w, r)
	}
}
