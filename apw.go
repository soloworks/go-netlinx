package apwfile

import (
	"errors"
	"os"
	"path/filepath"
)

// GatherFiles checks all referenced files exist, returns nil if all ok, otherwise returns array of filenames
func (apw *APWInfo) GatherFiles() error {

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
				// Depending on file path type, add this hash with full qualified name

				if filepath.IsAbs(file.FilePathName) {
					apw.Files[file.FilePathName] = file.Type
				} else {
					apw.Files[filepath.Join(apw.OriginPath, file.FilePathName)] = file.Type
				}

			}
		}
	}
	return nil
}

// checkAPWFiles verifies all files are present in current APWInfo references
func (apw *APWInfo) checkFiles() ([]string, error) {

	// Return Variable
	var badFiles []string

	// Cycle through all files in XML and check they exist
	for file := range apw.Files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			badFiles = append(badFiles, file)
		}
	}
	if len(badFiles) > 0 {
		return badFiles, errors.New("Files Missing")
	}
	return nil, nil
}
