package types

import (
	"time"
)

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

	Details *BugDetails
}

// Bugzilla is XML wrapper
type Bugzilla struct {
	Bug BugDetails `xml:"bug"`
}

// BugDetails contains additional information about bug. Filling this struct is
// expensive because it requires additional call to server for each bug.
type BugDetails struct {
	Version       string   `xml:"version"`
	TargetRelease string   `xml:"target_release"`
	Keywords      []string `xml:"keywords"`
}
