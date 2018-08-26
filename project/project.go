package project

import (
	"encoding/xml"
	"sort"

	"bitbucket.org/solo_works/apwfile/system"
)

// Project represetents an AMX project in an APW
type Project struct {
	XMLName       xml.Name         `xml:"Project"`
	Identifier    string           `xml:"Identifier"`
	Designer      string           `xml:"Designer,omitempty"`
	DealerID      string           `xml:"DealerID,omitempty"`
	SalesOrder    string           `xml:"SalesOrder,omitempty"`
	PurchaseOrder string           `xml:"PurchaseOrder,omitempty"`
	Comments      string           `xml:"Comments,omitempty"`
	Systems       []*system.System `xml:"System"`
}

// ByID implements sort.Interface for []Project based on
// the Identifier field.
type ByID []*Project

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].Identifier < a[j].Identifier }

// NewProject returns a pointer to a system
func NewProject(id string) *Project {
	return &Project{Identifier: id}
}

// FindSystem returns a pointer to a system
func (p *Project) FindSystem(id string) *system.System {
	for i, s := range p.Systems {
		if s.Identifier == id {
			return p.Systems[i]
		}
	}
	return nil
}

// AddSystem appends a project into a workspace and re-orders the projects
func (p *Project) AddSystem(s *system.System) {
	// Add this system to project
	p.Systems = append(p.Systems, s)
	// Sort the Systems
	sort.Sort(system.ByIdentifier(p.Systems))
}
