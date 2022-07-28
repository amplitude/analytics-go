package amplitude

import "time"

type EventOptions struct {
	UserID             string
	DeviceID           string
	Time               time.Time
	LocationLat        float64
	LocationLng        float64
	AppVersion         string
	VersionName        string
	Library            string
	Platform           string
	OSName             string
	OSVersion          string
	DeviceBrand        string
	DeviceManufacturer string
	DeviceModel        string
	Carrier            string
	Country            string
	Region             string
	City               string
	DMA                string
	IDFA               string
	IDFV               string
	ADID               string
	AndroidID          string
	Language           string
	IP                 string
	Price              float64
	Quantity           int
	Revenue            float64
	ProductID          string
	RevenueType        string
	EventID            int
	SessionID          int
	InsertID           string
	PartnerId          string
	Plan               Plan
}

type Event struct {
	EventType string
	EventOptions
	EventProperties map[string]interface{}
	UserProperties  map[IdentityOp]interface{}
	Groups          map[string][]string
	GroupProperties map[IdentityOp]interface{}
}

func (e Event) Clone() Event {
	return Event{
		EventType:       e.EventType,
		EventOptions:    e.EventOptions,
		EventProperties: cloneProperties(e.EventProperties),
		UserProperties:  cloneIdentiyProperties(e.UserProperties),
		Groups:          cloneGroups(e.Groups),
		GroupProperties: cloneIdentiyProperties(e.GroupProperties),
	}
}

// TODO: deep copy
func cloneProperties(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}

func cloneIdentiyProperties(properties map[IdentityOp]interface{}) map[IdentityOp]interface{} {
	result := make(map[IdentityOp]interface{}, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}

// TODO: deep copy
func cloneGroups(properties map[string][]string) map[string][]string {
	result := make(map[string][]string, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}
