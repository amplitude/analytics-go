package amplitude

import "time"

type EventOptions struct {
	UserID             string `json:"user_id"`
	DeviceID           string `json:"device_id"`
	Time               int64  `json:"time,omitempty"`
	InsertID           string `json:"insert_id,omitempty"`
	Library            string `json:"library,omitempty"`
	LocationLat        float64
	LocationLng        float64
	AppVersion         string
	VersionName        string
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
	PartnerId          string
	Plan               Plan
}

func (eo *EventOptions) setTime(time *time.Time) {
	eo.Time = time.UnixMilli()
}

type Event struct {
	EventType string `json:"event_type"`
	EventOptions
	EventProperties map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties  map[string]interface{} `json:"user_properties,omitempty"`
	Groups          map[string][]string    `json:"groups,omitempty"`
	GroupProperties map[string]interface{} `json:"group_properties,omitempty"`
}

func (e Event) Clone() Event {
	return Event{
		EventType:       e.EventType,
		EventOptions:    e.EventOptions,
		EventProperties: cloneProperties(e.EventProperties),
		UserProperties:  cloneProperties(e.UserProperties),
		Groups:          cloneGroups(e.Groups),
		GroupProperties: cloneProperties(e.GroupProperties),
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

// TODO: deep copy
func cloneGroups(properties map[string][]string) map[string][]string {
	result := make(map[string][]string, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}
