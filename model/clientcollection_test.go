package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestClientCollection(t *testing.T) {
	c := NewClientCollection()
	assert.Equal(t, c.Count(), 0, "new collections should not have items")

	c1 := &Client{ID: 1, Name: "First Client"}
	c.Add(c1)
	assert.Equal(t, c.Count(), 1, "Count() should correct value of clients in the collection")

	ct, ok := c.Get(1)
	assert.Assert(t, ok)
	assert.Equal(t, ct, c1)

	c2 := &Client{ID: 2, Name: "Second Client"}
	c.Add(c2)
	assert.Equal(t, c.Count(), 2, "Count() should correct value of clients in the collection")

	ct = c.Filter(func (client *Client) bool {
		return client.Name == "First Client"
	}).Take()
	assert.Equal(t, ct, c1, "Take() should return a client from the collection")

	var nc *Client
	ct = c.Filter(func (client *Client) bool {
		return client.Name == "Non-existing Client"
	}).Take()
	assert.Equal(t, ct, nc, "Take() should return nil when there are no clients in the collection")

	var cc int
	cc = c.Filter(func (client *Client) bool {
		return client.Name == "Non-existing Client"
	}).Count()
	assert.Equal(t, cc, 0, "Count() should correct value of clients in the collection")

	c3 := &Client{ID: 3, Name: "Third Client with an office", OfficeID: 333}
	c.Add(c3)
	ct = c.WithOffice(333).Take()
	assert.Equal(t, ct, c3)

	clients := c.All()
	assert.Equal(t, len(clients), 3)
}
