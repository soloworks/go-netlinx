package apw

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// APW holds details about an APW file
type APW struct {
	Identifier      string
	Name            string
	Filename        string
	OriginPath      string
	BuildPath       string
	FilesReferenced map[string]string
	FilesMissing    []string
	Workspace       *Workspace
}

// NewAPW loads or creates an APW object
sfunc NewAPW(filename string) *APW {

	// Create a new Structure
	var apw APW
	apw.FilesReferenced = make(map[string]string)

	// Populate the file details
	apw.Filename = filename
	apw.Name = filepath.Base(filename)
	apw.Identifier = strings.TrimSuffix(apw.Name, filepath.Ext(apw.Name))
	apw.OriginPath = filepath.Dir(filename)
	// Create a new empty Workspace
	apw.Workspace = &Workspace{}

	// Return
	return &apw
}

// LoadAPW loads or creates an APW object
func LoadAPW(filename string) (*APW, error) {

	// Make a new Object
	apw := NewAPW(filename)

	// Create or populate actual Workspace
	var err error
	apw.Workspace, err = Load(filename)
	if err == nil {
		// Gather File References
		apw.populateFileReferences()
	}

	// Return
	return apw, err
}

// GatherFiles checks all referenced files exist, returns nil if all ok, otherwise returns array of filenames
func (apw *APW) populateFileReferences() error {

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

// FindInFolder searches all subdirectories (recursivly option) for any
// .apw files and returns a list of AMXProjects
func FindInFolder(dir string, recursive bool) []*APW {

	// Create array to hold discovered projects
	var APWs []*APW

	// Check root folder and store APWs
	APWs = append(APWs, getAPWs(dir)...)

	// Do Sub Folder(s) if requested
	if recursive {
		files, _ := ioutil.ReadDir(dir)

		for _, subDir := range files {
			if subDir.IsDir() {
				APWs = append(APWs, FindInFolder(filepath.Join(dir, subDir.Name()), recursive)...)
			}
		}
	}

	// Return List
	return APWs
}

func getAPWs(dir string) []*APW {

	// Create array to hold discovered projects
	var APWs []*APW

	// Get all files in possible project folder
	files, _ := ioutil.ReadDir(dir)

	// Cycle through all and identify those with .apw files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".apw" {
			apw, err := LoadAPW(filepath.Join(dir, file.Name()))
			if err == nil {
				APWs = append(APWs, apw)
			}
		}
	}

	// Return List
	return APWs
}

// LoadWorkspace fetches the XML for this apw file
func (apw *APW) loadXML() error {
	var err error
	apw.Workspace, err = Load(apw.Filename)
	return err
}

// ExportWorkspace saves the XML to the designated destination folder
func (apw *APW) ExportWorkspace() error {
	Save(apw.Workspace, apw.Filename)
	return nil
}

// ExportArchive pulls all .apw files together into a zip in the target folder using the workspace name
func (apw *APW) ExportArchive(destDir string, buildID string) error {
	// Verify the APW file is all good before we do this
	if len(apw.FilesMissing) > 0 {
		var e bytes.Buffer
		e.WriteString(strconv.Itoa(len(apw.FilesMissing)))
		e.WriteString(" File")
		if len(apw.FilesMissing) > 1 {
			e.WriteString("s")
		}
		e.WriteString(" not found")
		return errors.New(e.String())
	}

	// Create a Zip writer with defered close
	var filename bytes.Buffer
	filename.WriteString(apw.Identifier)
	if buildID != "" {
		filename.WriteString("_" + buildID)
	}
	filename.WriteString(".zip")
	myZipFile, err := os.Create(filepath.Join(destDir, filename.String()))
	if err != nil {
		return err
	}
	defer myZipFile.Close()

	z := zip.NewWriter(myZipFile)
	defer z.Close()

	// Add each file to the Archive
	for file, fileType := range apw.FilesReferenced {

		// Open existing file
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)

		// Set file to correct folder based on file type
		header.Name = filepath.Join(FileFolder(fileType), header.Name)

		// Compress File
		header.Method = zip.Deflate

		writer, err := z.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}

	}

	// Save XML
	f, err := z.Create(apw.Identifier + ".apw")
	if err != nil {
		log.Fatal(err)
	}
	myXML, _ := apw.Workspace.Bytes()
	_, err = f.Write([]byte(myXML))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
