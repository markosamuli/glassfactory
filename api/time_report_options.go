package api

// TimeReportOptions represent the options that can be used when generating member time reports
type TimeReportOptions struct {
	clientID     int   // Client ID
	projectIDs   []int // Project IDs separated by comma
	officeID     int   // Office ID
	fetchRelated bool  // Fetch related data into the results
}

func (options *TimeReportOptions) apply(opts []TimeReportOption) {
	for _, o := range opts {
		o.apply(options)
	}
}

// TimeReportOption overrides behavior of TimeReport
type TimeReportOption interface {
	apply(*TimeReportOptions)
}

type timeReportOptionFunc func(*TimeReportOptions)

func (f timeReportOptionFunc) apply(o *TimeReportOptions) {
	f(o)
}

// WithClient returns time reports matching the given client
func WithClient(clientID int) TimeReportOption {
	return timeReportOptionFunc(func(o *TimeReportOptions) {
		o.clientID = clientID
	})
}

// WithProject returns time reports matching the given project
func WithProject(projectID int) TimeReportOption {
	return timeReportOptionFunc(func(o *TimeReportOptions) {
		o.projectIDs = append(o.projectIDs, projectID)
	})
}

// WithProjects returns time reports matching the given projects
func WithProjects(projectIDs []int) TimeReportOption {
	return timeReportOptionFunc(func(o *TimeReportOptions) {
		o.projectIDs = projectIDs
	})
}

// WithOffice returns time reports matching the given office
func WithOffice(officeID int) TimeReportOption {
	return timeReportOptionFunc(func(o *TimeReportOptions) {
		o.officeID = officeID
	})
}

// FetchRelated fetches the related data into the results
func FetchRelated() TimeReportOption {
	return timeReportOptionFunc(func(o *TimeReportOptions) {
		o.fetchRelated = true
	})
}

// NewTimeReportOptions returns TimeReportOptions with defaults
func NewTimeReportOptions(opts []TimeReportOption) *TimeReportOptions {
	options := &TimeReportOptions{
		fetchRelated: false, // Do not fetch related data by default
	}
	for _, o := range opts {
		o.apply(options)
	}
	return options
}
