package amplitude

type Identify struct {
	PropertiesSet []interface{}
	Properties    map[string]interface{}
}

// IsValid checks if to Identify object has Properties
// returns true if Identify object has Properties, otherwise returns false.
func (i Identify) IsValid() bool {
	return len(i.Properties) > 0
}

func (i Identify) Set(key string, value interface{}) {
}
