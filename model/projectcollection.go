package model


// ProjectCollection represents unique set of projects
type ProjectCollection struct {
	projects map[int]*Project
}

// NewProjectCollection is used for creating ProjectCollection
func NewProjectCollection() *ProjectCollection {
	c := &ProjectCollection{}
	c.projects = make(map[int]*Project, 0)
	return c
}

// Count returns number of projects in the collection
func (c *ProjectCollection) Count() int {
	return len(c.projects)
}

// Get returns a single project from the collection, if found
func (c *ProjectCollection) Get(id int) (*Project, bool) {
	project, ok := c.projects[id]
	return project, ok
}

// Add a project to the collection
func (c *ProjectCollection) Add(project *Project) {
	c.projects[project.ID] = project
}

// Take returns a project from the collection
func (c *ProjectCollection) Take() *Project {
	for _, project := range c.projects {
		return project
	}
	return nil
}

// All returns all projects from the collection
func (c *ProjectCollection) All() []*Project {
	projects := make([]*Project, 0)
	for _, project := range c.projects {
		projects = append(projects, project)
	}
	return projects
}

// Filter returns a new collection that contains projects matching the predicate
func (c *ProjectCollection) Filter(f func(*Project) bool) *ProjectCollection {
	fc := &ProjectCollection{}
	fc.projects = make(map[int]*Project, 0)
	for _, project := range c.projects {
		if f(project) {
			fc.projects[project.ID] = project
		}
	}
	return fc
}

// WithManager returns a new collection with projects matching the manager ID
func (c *ProjectCollection) WithManager(managerID int) *ProjectCollection {
	return c.Filter(func(m *Project) bool {
		return m.ManagerID == managerID
	})
}

// WithClient returns a new collection with projects matching the client ID
func (c *ProjectCollection) WithClient(clientID int) *ProjectCollection {
	return c.Filter(func(m *Project) bool {
		return m.ClientID == clientID
	})
}

// WithOffice returns a new collection with projects matching the office ID
func (c *ProjectCollection) WithOffice(officeID int) *ProjectCollection {
	return c.Filter(func(m *Project) bool {
		return m.OfficeID == officeID
	})
}

// WithJob returns a new collection with projects matching the job ID
func (c *ProjectCollection) WithJob(jobID string) *ProjectCollection {
	return c.Filter(func(m *Project) bool {
		return m.JobID == jobID
	})
}
