package workspace

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/soloworks/apwfile/project"
	"github.com/soloworks/apwfile/system"
)

// Load returns a struct containing the contents of the
// passed Project's .apw file
// Move the remove of global project to a seperate function
func Load(fn string) (*Workspace, error) {

	// Create new Workspace Object
	w := Workspace{}
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

	// Convert XML to Structure
	xml.Unmarshal(b, &w)

	// Sort Projects into order
	sort.Sort(project.ByID(w.Projects))

	// Sort Systems into order
	for _, p := range w.Projects {
		sort.Sort(system.ByIdentifier(p.Systems))
	}

	// Return Cleanly
	return &w, nil
}

// Save dumps the XML into the specified file
func Save(w *Workspace, fn string) error {
	// Create the dest directory
	os.MkdirAll(filepath.Dir(fn), os.ModePerm)
	// Dump XML to file
	b, _ := w.Bytes()
	err := ioutil.WriteFile(fn, b, 0644)
	return err
}

// Bytes converts structure to XML bytes
func (w *Workspace) Bytes() ([]byte, error) {

	// Set static values
	w.CurrentVersion = "4.0"
	// Create a ByteBuffer
	b := &bytes.Buffer{}
	// Add xml header constant
	b.WriteString(xml.Header)
	// Prep new encoder
	//output, err := xml.MarshalIndent(w, "", "")
	output, err := xml.Marshal(w)
	if err != nil {
		return nil, err
	}
	/*
		// Find each </ instance
		for _, s := range bytes.Split(output, []byte(">")) {
			b.Write(s)
			b.Write([]byte(">"))
			if bytes.Index(s, []byte("</")) != -1 {
				b.Write([]byte("\r\n"))
			}
		}
	*/
	b.Write(output)
	// Return Bytes
	return b.Bytes(), nil
}

// SetRelativeFilepaths sets all paths in the workspace to relative based on file type
func (w *Workspace) SetRelativeFilepaths() {
	// itterate over all elements
	for pi, p := range w.Projects {
		for si, s := range p.Systems {
			for fi, f := range s.Files {
				w.Projects[pi].Systems[si].Files[fi].FilePathName = filepath.Join(FileFolder(f.Type), filepath.Base(f.FilePathName))
			}
		}
	}
}

// SetAbsoluteFilepaths sets all relative paths in the workspace to absolute to match base directory provided
func (w *Workspace) SetAbsoluteFilepaths(path string) {
	// itterate over all elements
	for pi, p := range w.Projects {
		for si, s := range p.Systems {
			for fi, f := range s.Files {
				if !filepath.IsAbs(f.FilePathName) {
					w.Projects[pi].Systems[si].Files[fi].FilePathName = filepath.Join(path, f.FilePathName)
				}
			}
		}
	}
}

// FileFolder returns a string for a sub folder based on passed file type
func FileFolder(t string) string {
	switch t {
	case "TKN", "Source", "MasterSrc":
		{
			return "Source"
		}

	case "Include":
		{
			return "Includes"
		}

	case "IR", "AMX_IR_DB", "IRN_DB":
		{
			return "IR Files"
		}

	case "TP4", "TP5", "TPD", "KPD":
		{
			return "Interfaces"
		}

	case "XDD", "Module", "DUET", "TKO":
		{
			return "Modules"
		}
	default:
		{
			return "Other"
		}
	}
}
