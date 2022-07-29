package amplitude

type Identify struct {
	PropertiesSet []interface{}
	Properties    map[IdentityOp]interface{}
}

// isValid checks if to Identify object has Properties
// returns true if Identify object has Properties, otherwise returns false.
func (i *Identify) isValid() bool {
	return len(i.Properties) > 0
}

func (i *Identify) containsClearAllOperation() bool {
	for k := range i.Properties {
		if k == IdentityOpClearAll {
			return true
		}
	}

	return false
}

func (i *Identify) containsProperty(property interface{}) bool {
	for p := range i.PropertiesSet {
		if p == property {
			return true
		}
	}

	return false
}

func (i *Identify) containsOperation(op IdentityOp) bool {
	for operation := range i.Properties {
		if operation == op {
			return true
		}
	}

	return false
}

func (i *Identify) setUserProperty(op IdentityOp, property string, value interface{}) {
	if len(property) == 0 {
		// TO-DO: logger
	}

	if value == nil {
		// TO-DO: logger
	}

	if i.containsClearAllOperation() {
		// TO-DO: logger
	}

	if i.containsProperty(property) {
		// TO-DO: logger
	}

	if !i.containsOperation(op) {
		i.Properties[op] = make(map[string]interface{})
	}

	i.Properties[op] = map[string]interface{}{property: value}
	i.PropertiesSet = append(i.PropertiesSet, property)
}

// Set sets the value of a user property.
func (i *Identify) Set(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpSet, property, value)

	return i
}

// SetOnce sets the value of user property only once.
// Subsequent calls using SetOnce will be ignored.
func (i *Identify) SetOnce(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpSetOnce, property, value)

	return i
}

// Add increments a user property by some numerical value.
// If the user property does not have a value set yet,
// it will be initialized to 0 before being incremented.
func (i *Identify) Add(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpAdd, property, value)

	return i
}

// Prepend prepends a value or values to a user property array.
// If the user property does not have a value set yet,
// it will be initialized to an empty list before the new values are prepended.
func (i *Identify) Prepend(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpPrepend, property, value)

	return i
}

// Append appends a value or values to a user property array.
// If the user property does not have a value set yet,
// it will be initialized to an empty list before the new values are prepended.
func (i *Identify) Append(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpAppend, property, value)

	return i
}

// PreInsert pre-inserts a value or values to a user property,
// if it does not exist in the user property yet.
// Pre-insert means inserting the value(s) at the beginning of a given list.
// If the user property does not have a value set yet,
// it will be initialized to an empty list before the new values are pre-inserted.
// If the user property has an existing value, it will be no operation.
func (i *Identify) PreInsert(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpPreInsert, property, value)

	return i
}

// PostInsert post-inserts a value or values to a user property,
// if it does not exist in the user property yet.
// Post-insert means inserting the value(s) at the end of a given list.
// If the user property does not have a value set yet,
// it will be initialized to an empty list before the new values are post-inserted.
// If the user property has an existing value, it will be no operation.
func (i *Identify) PostInsert(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpPostInsert, property, value)

	return i
}

// Remove removes a value or values to a user property, if it exists in the user property.
// Remove means remove the existing value(s) from the given list.
// If the item does not exist in the user property, it will be no operation.
func (i *Identify) Remove(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpRemove, property, value)

	return i
}

// Unset removes the user property from the user profile.
func (i *Identify) Unset(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpUnset, property, value)

	return i
}

// ClearAll removes all user properties of this user.
func (i *Identify) ClearAll(property string, value interface{}) *Identify {
	i.setUserProperty(IdentityOpClearAll, property, value)

	return i
}
