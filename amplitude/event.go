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
	PartnerID          string
	Plan               Plan
}

type Event struct {
	EventType string
	EventOptions
	EventProperties map[string]interface{}
	UserProperties  map[IdentityOp]map[string]interface{}
	Groups          map[string][]string
	GroupProperties map[IdentityOp]map[string]interface{}
}

func (e Event) Clone() Event {
	return Event{
		EventType:       e.EventType,
		EventOptions:    e.EventOptions,
		EventProperties: cloneProperties(e.EventProperties),
		UserProperties:  cloneIdentityProperties(e.UserProperties),
		Groups:          cloneGroups(e.Groups),
		GroupProperties: cloneIdentityProperties(e.GroupProperties),
	}
}

func cloneProperties(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range properties {
		vm, ok := v.(map[string]interface{})
		if ok {
			result[k] = cloneProperties(vm)
		} else {
			result[k] = v
		}
	}

	return result
}

func cloneIdentityProperties(properties map[IdentityOp]map[string]interface{}) map[IdentityOp]map[string]interface{} {
	result := make(map[IdentityOp]map[string]interface{})

	for operation, p := range properties {
		result[operation] = cloneProperties(p)
	}

	return result
}

func cloneGroups(properties map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for k, v := range properties {
		result[k] = make([]string, len(v))
		copy(result[k], v)
	}

	return result
}
