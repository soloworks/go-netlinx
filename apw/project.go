package apw

import (
	"encoding/xml"
	"sort"
)

// Project represetents an AMX project in an APW
type Project struct {
	XMLName       xml.Name  `xml:"Project"`
	Identifier    string    `xml:"Identifier"`
	Designer      string    `xml:"Designer,omitempty"`
	DealerID      string    `xml:"DealerID,omitempty"`
	SalesOrder    string    `xml:"SalesOrder,omitempty"`
	PurchaseOrder string    `xml:"PurchaseOrder,omitempty"`
	Comments      string    `xml:"Comments,omitempty"`
	Systems       []*System `xml:"System"`
}

// ByProjectID implements sort.Interface for []Project based on
// the Identifier field.
type ByProjectID []*Project

func (a ByProjectID) Len() int           { return len(a) }
func (a ByProjectID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProjectID) Less(i, j int) bool { return a[i].Identifier < a[j].Identifier }

// NewProject returns a pointer to a system
func NewProject(id string) *Project {
	return &Project{Identifier: id}
}

// FindSystem returns a pointer to a system
func (p *Project) FindSystem(id string) *System {
	for i, s := range p.Systems {
		if s.Identifier == id {
			return p.Systems[i]
		}
	}
	return nil
}

// AddSystem appends a project into a workspace and re-orders the projects
func (p *Project) AddSystem(s *System) {
	// Add this system to project
	p.Systems = append(p.Systems, s)
	// Sort the Systems
	sort.Sort(BySystemID(p.Systems))
}
