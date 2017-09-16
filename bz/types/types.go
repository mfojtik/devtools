package types

import "time"

// Buglist holds a list of bugs.
type Buglist struct {
	Items []Bug
}

// Bug represents a signle bug details as parsed from CSV.
type Bug struct {
	BugID      string
	Product    string
	Component  string
	Assignee   string
	Status     string
	Resolution string
	Summary    string
	Changed    time.Time
}
