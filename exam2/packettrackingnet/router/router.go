package router

const (
	SHIPMENTS             = "/shipments"
	PACKETS               = "/packets"
	PACKETS_RECEIVED      = PACKETS + "/received"
	LOCATIONS             = "/locations"
	TRACKING              = "/tracking"
	SERVICES              = "/services"
	SERVICES_NAMES        = SERVICES + "/names"
	SHIPMENTS_SENDERS     = SHIPMENTS + "/senders"
	SHIPMENTS_RECEIVERS   = SHIPMENTS + "/receivers"
	SHIPMENTS_BULK_CREATE = SHIPMENTS + "/bulk" + "/create"
)
