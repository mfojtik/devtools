package godep

import (
	"encoding/json"
	"fmt"
	"os"
)

// Godeps describes what a package needs to be rebuilt reproducibly.
// It's the same information stored in file Godeps.
type Godeps struct {
	ImportPath   string
	GoVersion    string
	GodepVersion string
	Packages     []string `json:",omitempty"` // Arguments to save, if any.
	Deps         []Dependency
	isOldFile    bool
}

// A Dependency is a specific revision of a package.
type Dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"` // Description of commit, if present.
	Rev        string // VCS-specific commit ID.
}

// LoadGodepsFile loads the godep file
func LoadGodepsFile(path string) (Godeps, error) {
	var g Godeps
	f, err := os.Open(path)
	if err != nil {
		return g, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&g)
	if err != nil {
		err = fmt.Errorf("Unable to parse %s: %s", path, err.Error())
	}
	return g, err
}
