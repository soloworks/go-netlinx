package apw

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// Workspace represents the APW XML file structure
type Workspace struct {
	XMLName        xml.Name   `xml:"Workspace"`
	Identifier     string     `xml:"Identifier"`
	CreateVersion  string     `xml:"CreateVersion"`
	PJSFile        string     `xml:"PJS_File,omitempty"`
	PJSConvertDate string     `xml:"PJS_ConvertDate,omitempty"`
	PJSCreateDate  string     `xml:"PJS_CreateDate,omitempty"`
	Comments       string     `xml:"Comments,omitempty"`
	Projects       []*Project `xml:"Project"`
	CurrentVersion string     `xml:"CurrentVersion,attr"`
}

// NewWorkspace returns a new workspace instance with
// default fields already populated
func NewWorkspace(identifier string) Workspace {
	return Workspace{
		Identifier:     identifier,
		CurrentVersion: "4.0",
		CreateVersion:  "4.0",
	}
}

// FindProject returns a pointer to the required project
func (w *Workspace) FindProject(id string) *Project {
	for i, p := range w.Projects {
		if p.Identifier == id {
			return w.Projects[i]
		}
	}
	return nil
}

// AddProject appends a project into a workspace and re-orders the projects
func (w *Workspace) AddProject(p *Project) {
	// Add this system to project
	w.Projects = append(w.Projects, p)
	// Sort the Systems
	sort.Sort(ByProjectID(w.Projects))
}

type packType int

const (
	ptRelease packType = iota
	ptHandover
	ptFull
)

// ConvertToHandover Converts all files to Handover Type
func (w *Workspace) ConvertToHandover() {
	w.convert(ptHandover)
}

// ConvertToRelease Converts all files to Handover Type
func (w *Workspace) ConvertToRelease() {
	w.convert(ptRelease)
}

// Convert sets all files to type based on package passed
func (w *Workspace) convert(pt packType) {
	// itterate over all Projects
	for pi, p := range w.Projects {
		// Itterate over all Systems
		for si, s := range p.Systems {
			// Create new File Stack for this system
			var files []*File
			// Itterate over all Files
			for _, f := range s.Files {
				// Select the extention
				switch filepath.Ext(f.FilePathName) {
				// For .axs files (Main Source Files)
				case ".axs":
					switch f.Type {
					case "Source", "MasterSrc":
						switch pt {
						case ptRelease:
							// Swap the extension to set file to compiled source code
							f.ChangeExtension("tkn")
							if f.DeviceMaps == nil {
								f.AddDeviceMap(NewDeviceMap("Custom [0:1:0]", "Custom [0:1:0]"))
							}
						}

					case "Module":
						// Assign file in correct places
						switch pt {
						case ptRelease:
							f = nil
						}
						switch pt {
						case ptHandover:
							// Swap the extension to set file to compiled module
							f.ChangeExtension("tko")
						}
					}
				case ".axi":
					switch f.Type {
					case "Include":
						switch pt {
						case ptRelease:
							f = nil
						}
					}
				}
				if f != nil {
					files = append(files, f)
				}
			}
			w.Projects[pi].Systems[si].Files = files
		}
	}
}

// Load returns a struct containing the contents of the
// passed Project's .apw file
// Move the remove of global project to a seperate function
// Returns empty workspace if file doesn't exist
func Load(fn string) (*Workspace, error) {

	// Create new Workspace Object
	w := Workspace{}
	// Open to .apw file
	f, err := os.Open(fn)
	if err != nil {
		return &w, err
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
	sort.Sort(ByProjectID(w.Projects))

	// Sort Systems into order
	for _, p := range w.Projects {
		sort.Sort(BySystemID(p.Systems))
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
