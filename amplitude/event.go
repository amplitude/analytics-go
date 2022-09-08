package amplitude

import "time"

type EventOptions struct {
	UserID             string  `json:"user_id,omitempty"`
	DeviceID           string  `json:"device_id,omitempty"`
	Time               int64   `json:"time,omitempty"`
	InsertID           string  `json:"insert_id,omitempty"`
	Library            string  `json:"library,omitempty"`
	LocationLat        float64 `json:"location_lat,omitempty"`
	LocationLng        float64 `json:"location_lng,omitempty"`
	AppVersion         string  `json:"app_version,omitempty"`
	VersionName        string  `json:"version_name,omitempty"`
	Platform           string  `json:"platform,omitempty"`
	OSName             string  `json:"os_name,omitempty"`
	OSVersion          string  `json:"os_version,omitempty"`
	DeviceBrand        string  `json:"device_brand,omitempty"`
	DeviceManufacturer string  `json:"device_manufacturer,omitempty"`
	DeviceModel        string  `json:"device_model,omitempty"`
	Carrier            string  `json:"carrier,omitempty"`
	Country            string  `json:"country,omitempty"`
	Region             string  `json:"region,omitempty"`
	City               string  `json:"city,omitempty"`
	DMA                string  `json:"dma,omitempty"`
	IDFA               string  `json:"idfa,omitempty"`
	IDFV               string  `json:"idfv,omitempty"`
	ADID               string  `json:"adid,omitempty"`
	AndroidID          string  `json:"android_id,omitempty"`
	Language           string  `json:"language,omitempty"`
	IP                 string  `json:"ip,omitempty"`
	Price              float64 `json:"price,omitempty"`
	Quantity           int     `json:"quantity,omitempty"`
	Revenue            float64 `json:"revenue,omitempty"`
	ProductID          string  `json:"productId,omitempty"`
	RevenueType        string  `json:"revenueType,omitempty"`
	EventID            int     `json:"event_id,omitempty"`
	SessionID          int     `json:"session_id,omitempty"`
	PartnerID          string  `json:"partner_id,omitempty"`
	Plan               Plan    `json:"plan,omitempty"`
}

func (eo *EventOptions) SetTime(time time.Time) {
	eo.Time = time.UnixMilli()
}

type Event struct {
	EventType string `json:"event_type"`
	EventOptions
	EventProperties map[string]interface{}                `json:"event_properties,omitempty"`
	UserProperties  map[IdentityOp]map[string]interface{} `json:"user_properties,omitempty"`
	Groups          map[string][]string                   `json:"groups,omitempty"`
	GroupProperties map[IdentityOp]map[string]interface{} `json:"group_properties,omitempty"`
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
	if properties == nil {
		return nil
	}

	clone := make(map[string]interface{}, len(properties))

	for k, v := range properties {
		clone[k] = cloneUnknown(v)
	}

	return clone
}

func cloneIdentityProperties(properties map[IdentityOp]map[string]interface{}) map[IdentityOp]map[string]interface{} {
	if properties == nil {
		return nil
	}

	clone := make(map[IdentityOp]map[string]interface{})

	for operation, p := range properties {
		clone[operation] = cloneProperties(p)
	}

	return clone
}

func cloneGroups(properties map[string][]string) map[string][]string {
	if properties == nil {
		return nil
	}

	clone := make(map[string][]string, len(properties))
	for k, v := range properties {
		clone[k] = make([]string, len(v))
		copy(clone[k], v)
	}

	return clone
}

func cloneIntegers(values []int) []int {
	if values == nil {
		return nil
	}

	clone := make([]int, len(values))
	copy(clone, values)

	return clone
}

func cloneFloats(values []float64) []float64 {
	if values == nil {
		return nil
	}

	clone := make([]float64, len(values))
	copy(clone, values)

	return clone
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}

	clone := make([]string, len(values))
	copy(clone, values)

	return clone
}

func cloneBooleans(values []bool) []bool {
	if values == nil {
		return nil
	}

	clone := make([]bool, len(values))
	copy(clone, values)

	return clone
}

func cloneUnknowns(values []interface{}) []interface{} {
	if values == nil {
		return nil
	}

	clone := make([]interface{}, len(values))
	for i, value := range values {
		clone[i] = cloneUnknown(value)
	}

	return clone
}

func cloneUnknown(value interface{}) interface{} {
	switch value := value.(type) {
	case []int:
		return cloneIntegers(value)
	case []float64:
		return cloneFloats(value)
	case []string:
		return cloneStrings(value)
	case []bool:
		return cloneBooleans(value)
	case []interface{}:
		return cloneUnknowns(value)
	case map[string]interface{}:
		clone := make(map[string]interface{}, len(value))
		for k, v := range clone {
			clone[k] = cloneUnknown(v)
		}

		return clone
	default:
		return value
	}
}
