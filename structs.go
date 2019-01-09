package apwfile

import (
	"path/filepath"
	"strings"

	"github.com/soloworks/go-netlinx-apwfile/workspace"
)

// APWInfo holds details about an APW file
type APWInfo struct {
	Identifier      string
	Name            string
	Filename        string
	OriginPath      string
	BuildPath       string
	FilesReferenced map[string]string
	FilesMissing    []string
	Workspace       *workspace.Workspace
}

// NewAPWInfo loads or creates an APWInfo object
func NewAPWInfo(filename string) *APWInfo {

	// Create a new Structure
	var apw APWInfo
	apw.FilesReferenced = make(map[string]string)

	// Populate the file details
	apw.Filename = filename
	apw.Name = filepath.Base(filename)
	apw.Identifier = strings.TrimSuffix(apw.Name, filepath.Ext(apw.Name))
	apw.OriginPath = filepath.Dir(filename)
	// Create a new empty Workspace
	apw.Workspace = &workspace.Workspace{}

	// Return
	return &apw
}

// LoadAPWInfo loads or creates an APWInfo object
func LoadAPWInfo(filename string) (*APWInfo, error) {

	// Make a new Object
	apw := NewAPWInfo(filename)

	// Create or populate actual Workspace
	var err error
	apw.Workspace, err = workspace.Load(filename)
	if err == nil {
		// Gather File References
		apw.populateFileReferences()
	}

	// Return
	return apw, err
}
