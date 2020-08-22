package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestProjectCollection(t *testing.T) {
	c := NewProjectCollection()
	assert.Equal(t, c.Count(), 0)

	p1 := &Project{ID: 1, Name: "First Project", OfficeID: 111}
	c.Add(p1)
	assert.Equal(t, c.Count(), 1)

	var p *Project
	var ok bool
	p, ok = c.Get(1)
	assert.Assert(t, ok)
	assert.Equal(t, p, p1)

	p2 := &Project{ID: 2, Name: "Second Project", ClientID: 222}
	c.Add(p2)
	assert.Equal(t, c.Count(), 2)

	p = c.Filter(func(member *Project) bool {
		return member.Name == "First Project"
	}).Take()
	assert.Equal(t, p, p1)

	var pc int
	pc = c.Filter(func(member *Project) bool {
		return member.Name == "Non-existing Project"
	}).Count()
	assert.Equal(t, pc, 0)

	p = c.Filter(func(member *Project) bool {
		return member.Name == "Non-existing Project"
	}).Take()
	var np *Project
	assert.Equal(t, p, np)

	p = c.WithOffice(111).Take()
	assert.Equal(t, p, p1)

	p = c.WithClient(222).Take()
	assert.Equal(t, p, p2)

	p3 := &Project{ID: 3, Name: "Third Project with Job ID", JobID: "JOB333"}
	c.Add(p3)
	p = c.WithJob("JOB333").Take()
	assert.Equal(t, p, p3)

	p4 := &Project{ID: 4, Name: "Fourth Project with a manager", ManagerID: 444}
	c.Add(p4)
	p = c.WithManager(444).Take()
	assert.Equal(t, p, p4)

	projects := c.All()
	assert.Equal(t, len(projects), 4)
}
