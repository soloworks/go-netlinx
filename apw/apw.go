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
	"sort"
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
func NewAPW(fn string, xml []byte) (*APW, error) {

	// Create a new Structure
	var apw APW
	apw.FilesReferenced = make(map[string]string)

	// Populate the file details
	apw.Filename = fn
	apw.Name = filepath.Base(fn)
	apw.Identifier = strings.TrimSuffix(apw.Name, filepath.Ext(apw.Name))
	apw.OriginPath = filepath.Dir(fn)
	// Create a new empty Workspace
	apw.Workspace = &Workspace{}
	// Populate and Process if xml is present
	if xml != nil {
		err := apw.Workspace.FromBytes(xml)
		if err != nil {
			return nil, err
		}

		// Gather File References
		apw.populateFileReferences()
		// Sort Projects into order
		sort.Sort(ByProjectID(apw.Workspace.Projects))

		// Sort Systems into order
		for _, p := range apw.Workspace.Projects {
			sort.Sort(BySystemID(p.Systems))
		}

	}

	// Return
	return &apw, nil
}

// LoadAPW loads an APW from filename, the passes to ParseAPW for Return
func LoadAPW(fn string) (*APW, error) {

	// Open to .apw file
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read contents of file
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Make a new Object
	apw, err := NewAPW(fn, b)
	if err != nil {
		return nil, err
	}
	// Return
	return apw, err
}

// writeFile outputs the workspace XML into the specified file
func writeFile(w *Workspace, fn string) error {
	// Create the dest directory
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	// Dump XML to file
	b, _ := w.ToXML()
	err := ioutil.WriteFile(fn, b, 0644)
	return err
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

// FindAPWs searches all subdirectories (recursivly option) for any
// .apw files and returns a list of AMXProjects
func FindAPWs(sourceDir string, recursive bool) []*APW {

	// Create array to hold discovered projects
	var APWs []*APW

	// Check root folder and store APWs
	APWs = append(APWs, gatherAPWs(sourceDir)...)

	// Do Sub Folder(s) if requested
	if recursive {
		files, _ := ioutil.ReadDir(sourceDir)

		for _, subDir := range files {
			if subDir.IsDir() {
				APWs = append(APWs, FindAPWs(filepath.Join(sourceDir, subDir.Name()), recursive)...)
			}
		}
	}

	// Return List
	return APWs
}

func gatherAPWs(dir string) []*APW {

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

// ExportAPW saves the XML to the designated destination folder
func (apw *APW) ExportAPW() error {
	writeFile(apw.Workspace, apw.Filename)
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
	b, _ := apw.Workspace.ToXML()
	_, err = f.Write([]byte(b))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
