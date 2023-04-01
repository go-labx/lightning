package lightning

// ContextData is a map[string]interface{} that can be used to store data in the context of a request.
type ContextData map[string]interface{}

// Get retrieves the value associated with the given key from the ContextData.
func (c ContextData) Get(key string) interface{} {
	return c[key]
}

// Set sets the value associated with the given key in the ContextData.
func (c ContextData) Set(key string, value interface{}) {
	c[key] = value
}

// Del deletes the value associated with the given key from the ContextData.
func (c ContextData) Del(key string) {
	delete(c, key)
}
