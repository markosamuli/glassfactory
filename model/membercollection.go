package model

// MemberCollection represents unique set of members
type MemberCollection struct {
	members map[int]*Member
}

// NewMemberCollection is used for creating MemberCollection
func NewMemberCollection() *MemberCollection {
	c := &MemberCollection{}
	c.members = make(map[int]*Member, 0)
	return c
}

// Count returns number of members in the collection
func (c *MemberCollection) Count() int {
	return len(c.members)
}

// Get returns a single member from the collection, if found
func (c *MemberCollection) Get(id int) (*Member, bool) {
	member, ok := c.members[id]
	return member, ok
}

// Add a member to the collection
func (c *MemberCollection) Add(member *Member) {
	c.members[member.ID] = member
}

// Take returns a member from the collection
func (c *MemberCollection) Take() *Member {
	for _, member := range c.members {
		return member
	}
	return nil
}

// All returns all members from the collection
func (c *MemberCollection) All() []*Member {
	members := make([]*Member, 0)
	for _, member := range c.members {
		members = append(members, member)
	}
	return members
}

// Filter returns a new collection that contains members matching the predicate
func (c *MemberCollection) Filter(f func(*Member) bool) *MemberCollection {
	fc := &MemberCollection{}
	fc.members = make(map[int]*Member, 0)
	for _, member := range c.members {
		if f(member) {
			fc.members[member.ID] = member
		}
	}
	return fc
}

// WithEmail returns a new collection with members matching the email
func (c *MemberCollection) WithEmail(email string) *MemberCollection {
	return c.Filter(func(m *Member) bool {
		return m.Email == email
	})
}

// WithRole returns a new collection with members matching the role ID
func (c *MemberCollection) WithRole(roleID int) *MemberCollection {
	return c.Filter(func(m *Member) bool {
		return m.RoleID == roleID
	})
}

// WithOffice returns a new collection with members matching the office ID
func (c *MemberCollection) WithOffice(officeID int) *MemberCollection {
	return c.Filter(func(m *Member) bool {
		return m.OfficeID == officeID
	})
}
