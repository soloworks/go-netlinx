package apwfile

import (
	"errors"
	"os"
	"path/filepath"
)

// GatherFiles checks all referenced files exist, returns nil if all ok, otherwise returns array of filenames
func (apw *APWInfo) populateFileReferences() error {

	// Return Error if no XML is present
	if &apw.Workspace == nil {
		return errors.New("XML Empty")
	}

	// Cycle through Projects
	for _, project := range apw.Workspace.Projects {
		// Cycle through Systems
		for _, system := range project.Systems {
			// Cycle through Files
			for _, file := range system.Files {
				// Depending on file path type (Absolute or Relative), add this hash with full qualified name
				if filepath.IsAbs(file.FilePathName) {
					apw.FilesReferenced[file.FilePathName] = file.Type
				} else {
					apw.FilesReferenced[filepath.Join(apw.OriginPath, file.FilePathName)] = file.Type
				}

			}
		}
	}

	// Cycle through all files in XML and check they exist
	for file := range apw.FilesReferenced {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			apw.FilesMissing = append(apw.FilesMissing, file)
		}
	}

	return nil
}
