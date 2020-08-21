package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestMemberCollection(t *testing.T) {
	c := NewMemberCollection()
	assert.Equal(t, c.Count(), 0)

	m1 := &Member{ID: 1, Name: "First User", Email: "first@email.com"}
	c.Add(m1)
	assert.Equal(t, c.Count(), 1)

	var (
		m *Member
		ok bool
	)
	m, ok = c.Get(1)
	assert.Assert(t, ok)
	assert.Equal(t, m, m1)

	m2 := &Member{ID: 2, Name: "Second User", Email: "second@email.com"}
	c.Add(m2)
	assert.Equal(t, c.Count(), 2)

	m = c.Filter(func(member *Member) bool {
		return member.Name == "First User"
	}).Take()
	assert.Equal(t, m, m1)

	var mc int
	mc = c.Filter(func(member *Member) bool {
		return member.Name == "Non-existing User"
	}).Count()
	assert.Equal(t, mc, 0)

	var nm *Member
	m = c.Filter(func(member *Member) bool {
		return member.Name == "Non-existing User"
	}).Take()
	assert.Equal(t, m, nm)

	m = c.WithEmail("second@email.com").Take()
	assert.Equal(t, m, m2)

	m3 := &Member{ID: 3, Name: "Third User with a role", RoleID: 333}
	c.Add(m3)
	m = c.WithRole(333).Take()
	assert.Equal(t, m, m3)

	m4 := &Member{ID: 4, Name: "Forth User with an office", OfficeID: 444}
	c.Add(m4)
	m = c.WithOffice(444).Take()
	assert.Equal(t, m, m4)

	m5 := &Member{ID: 2, Name: "Second User with new details", Email: "second.changed@email.com"}
	c.Add(m5)
	m, ok = c.Get(2)
	assert.Assert(t, ok)
	assert.Equal(t, m, m5)

	members := c.All()
	assert.Equal(t, len(members), 4)
}