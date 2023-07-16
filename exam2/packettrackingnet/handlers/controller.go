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
type System struct{}

func (sh *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.SYSTEM_DB_TRUNCATE) && r.Method == http.MethodGet:
		truncateData(w)
	default:
		endpointNotFound(w)
	}
}

func (sh *ShipmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.SHIPMENTS_SENDERS) && r.Method == http.MethodPost:
		postSender(w, r)
	case strings.EqualFold(path, router.SHIPMENTS_RECEIVERS) && r.Method == http.MethodPost:
		postReceiver(w, r)
	case strings.EqualFold(path, router.SHIPMENTS) && r.Method == http.MethodPost:
		postShipment(w, r)
	case strings.EqualFold(path, router.SHIPMENTS) && r.Method == http.MethodGet:
		getAllShipments(w, r)
	case strings.EqualFold(path, router.SHIPMENTS) && r.Method == http.MethodPut:
		updateShipmentCheckpoint(w, r)
	case strings.EqualFold(path, router.SHIPMENTS_BULK_CREATE) && r.Method == http.MethodPost:
		bulkCreateShipments(w, r)
	case strings.EqualFold(path, router.SHIPMENTS_DOWNLOAD) && r.Method == http.MethodGet:
		downloadShipmentCSV(w)
	default:
		endpointNotFound(w)
	}
}

func (lh *LocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.EqualFold(path, router.LOCATIONS) && r.Method == http.MethodPost:
		postLocation(w, r)
	case strings.EqualFold(path, router.LOCATIONS) && r.Method == http.MethodGet:
		getLocations(w, r)
	case strings.EqualFold(path, router.LOCATIONS) && r.Method == http.MethodPut:
		updateLocationAddressByName(w, r)
	default:
		endpointNotFound(w)
	}
}

func (ph *PacketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()
	switch {
	case strings.EqualFold(path, router.PACKETS) && r.Method == http.MethodPost:
		postPacket(w, r)
	case strings.Contains(path, "locationName") && r.Method == http.MethodGet:
		getAllPacketByLocationName(w, r)
	case strings.EqualFold(path, router.PACKETS_RECEIVED) && r.Method == http.MethodGet:
		getAllReceivedPackets(w, r)
	default:
		endpointNotFound(w)
	}
}

func (th *TrackingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.Contains(path, "trackingId") && r.Method == http.MethodGet:
		getShipmentById(w, r)
	default:
		endpointNotFound(w)
	}
}

func (svh *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.String()

	switch {
	case strings.Contains(path, "serviceName") && r.Method == http.MethodGet:
		getServiceByName(w, r)
	case strings.EqualFold(path, router.SERVICES) && r.Method == http.MethodPost:
		postService(w, r)
	case strings.EqualFold(path, router.SERVICES) && r.Method == http.MethodGet:
		getServices(w, r)
	case strings.EqualFold(path, router.SERVICES_NAMES) && r.Method == http.MethodGet:
		getAllServiceNames(w, r)
	default:
		endpointNotFound(w)
	}
}
