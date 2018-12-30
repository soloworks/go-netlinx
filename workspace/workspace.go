package workspace

import (
	"encoding/xml"
	"path/filepath"
	"sort"

	"github.com/soloworks/go-netlinx-apwfile/devicemap"
	"github.com/soloworks/go-netlinx-apwfile/file"
	"github.com/soloworks/go-netlinx-apwfile/project"
)

// Workspace represents the APW XML file structure
type Workspace struct {
	XMLName        xml.Name           `xml:"Workspace"`
	Identifier     string             `xml:"Identifier"`
	CreateVersion  string             `xml:"CreateVersion"`
	PJSFile        string             `xml:"PJS_File,omitempty"`
	PJSConvertDate string             `xml:"PJS_ConvertDate,omitempty"`
	PJSCreateDate  string             `xml:"PJS_CreateDate,omitempty"`
	Comments       string             `xml:"Comments,omitempty"`
	Projects       []*project.Project `xml:"Project"`
	CurrentVersion string             `xml:"CurrentVersion,attr"`
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
func (w *Workspace) FindProject(id string) *project.Project {
	for i, p := range w.Projects {
		if p.Identifier == id {
			return w.Projects[i]
		}
	}
	return nil
}

// AddProject appends a project into a workspace and re-orders the projects
func (w *Workspace) AddProject(p *project.Project) {
	// Add this system to project
	w.Projects = append(w.Projects, p)
	// Sort the Systems
	sort.Sort(project.ByID(w.Projects))
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
			var files []*file.File
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
								f.AddDeviceMap(devicemap.NewDeviceMap("Custom [0:1:0]", "Custom [0:1:0]"))
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
