package apwfile

import (
	"path/filepath"
	"strings"

	"github.com/soloworks/apwfile/workspace"
)

// APWInfo holds details about an APW file
type APWInfo struct {
	Identifier string
	Name       string
	Filename   string
	OriginPath string
	BuildPath  string
	Files      map[string]string
	Workspace  *workspace.Workspace
}

// NewAPWInfo gets a list of unique files referenced in a Workspace
func NewAPWInfo(filename string) *APWInfo {

	// Create a new Structure
	var apw APWInfo

	// Populate the file details
	apw.Filename = filename
	apw.Name = filepath.Base(filename)
	apw.Identifier = strings.TrimSuffix(apw.Name, filepath.Ext(apw.Name))
	apw.OriginPath = filepath.Dir(filename)
	apw.Files = make(map[string]string)

	// Create a new empty Workspace
	apw.Workspace = &workspace.Workspace{}

	// Return
	return &apw
}
