package model

// ClientCollection represents unique set of clients
type ClientCollection struct {
	clients map[int]*Client
}

// NewClientCollection is used for creating ClientCollection
func NewClientCollection() *ClientCollection {
	c := &ClientCollection{}
	c.clients = make(map[int]*Client, 0)
	return c
}

// Count returns number of clients in the collection
func (c *ClientCollection) Count() int {
	return len(c.clients)
}

// Get returns a single client from the collection, if found
func (c *ClientCollection) Get(id int) (*Client, bool) {
	client, ok := c.clients[id]
	return client, ok
}

// Add a client to the collection
func (c *ClientCollection) Add(client *Client) {
	c.clients[client.ID] = client
}

// Take returns a client from the collection
func (c *ClientCollection) Take() *Client {
	for _, client := range c.clients {
		return client
	}
	return nil
}

// All returns all clients from the collection
func (c *ClientCollection) All() []*Client {
	clients := make([]*Client, 0)
	for _, client := range c.clients {
		clients = append(clients, client)
	}
	return clients
}

// Filter returns a new collection that contains clients matching the predicate
func (c *ClientCollection) Filter(f func(*Client) bool) *ClientCollection {
	fc := &ClientCollection{}
	fc.clients = make(map[int]*Client, 0)
	for _, client := range c.clients {
		if f(client) {
			fc.clients[client.ID] = client
		}
	}
	return fc
}

// WithOffice returns a new collection with clients matching the office ID
func (c *ClientCollection) WithOffice(officeID int) *ClientCollection {
	return c.Filter(func(m *Client) bool {
		return m.OfficeID == officeID
	})
}