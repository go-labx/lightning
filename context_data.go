package lightning

// contextData is a map[string]interface{} that can be used to store data in the context of a originReq.
type contextData map[string]interface{}

// get retrieves the value associated with the given key from the contextData.
func (c contextData) get(key string) interface{} {
	return c[key]
}

// set sets the value associated with the given key in the contextData.
func (c contextData) set(key string, value interface{}) {
	c[key] = value
}

// del deletes the value associated with the given key from the contextData.
func (c contextData) del(key string) {
	delete(c, key)
}
