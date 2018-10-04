package file

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/soloworks/apwfile/devicemap"
	"github.com/soloworks/apwfile/irdb"
)

// A Type specifies a month of the year (January = 1, ...).
type Type int

// File Types for use outside this module
const (
	TypeSource Type = iota
	TypeMasterSrc
	TypeInclude
	TypeModule
	TypeAXB
	TypeIR
	TypeTPD
	TypeTP4
	TypeTP5
	TypeKPD
	TypeTKO
	TypeIRDB
	TypeIRNDB
	TypeOther
	TypeDuet
	TypeTOK
	TypeTKN
	TypeKPB
	TypeXDD
)

// Types for use outside this module
var types = [...]string{
	"Source",
	"MasterSrc",
	"Include",
	"Module",
	"AXB",
	"IR",
	"TPD",
	"TP4",
	"TP5",
	"KPD",
	"TKO",
	"IRDB",
	"IRNDB",
	"Other",
	"Duet",
	"TOK",
	"TKN",
	"KPB",
	"XDD",
}

// String returns the English name of the Type
func (t Type) String() string { return types[t] }

// CompileType specifies a file compilation type
type CompileType int

// File Types for use outside this module
const (
	CompileTypeNone CompileType = iota
	CompileTypeNetlinx
	CompileTypeAxcess
)

// Types for use outside this module
var compileTypes = [...]string{
	"None",
	"Netlinx",
	"Axcess",
}

// String returns the English name of the Type
func (ct CompileType) String() string { return compileTypes[ct] }

// File represetents an AMX project in an APW
type File struct {
	XMLName         xml.Name               `xml:"File"`
	Identifier      string                 `xml:"Identifier"`
	FilePathName    string                 `xml:"FilePathName"`
	Comments        string                 `xml:"Comments,omitempty"`
	MasterDirectory string                 `xml:"MasterDirectory,omitempty"`
	DeviceMaps      []*devicemap.DeviceMap `xml:"DeviceMap"`
	IRDBs           []*irdb.IRDB           `xml:"IRDB"`
	Type            string                 `xml:"Type,attr"`
	CompileType     string                 `xml:"CompileType,attr"`
}

// NewFile returns a new project instance with
// default fields already populated
func NewFile(f string, t Type, c CompileType) *File {
	return &File{
		Identifier:   filepath.Base(f),
		FilePathName: f,
		Type:         t.String(),
		CompileType:  c.String(),
	}
}

// AddDeviceMap adds a file to a system
func (f *File) AddDeviceMap(d *devicemap.DeviceMap) {
	f.DeviceMaps = append(f.DeviceMaps, d)
}

// ChangeExtension will replace the existing file extension for the passed value
func (f *File) ChangeExtension(ext string) {
	// Create a buffer to build this new file name up
	var b bytes.Buffer
	// Get the existing file without extension
	b.WriteString(strings.TrimSuffix(f.FilePathName, filepath.Ext(f.FilePathName)))
	// Add the new Extension
	b.WriteString(".")
	b.WriteString(ext)
	// Replace the existing filename with the new
	f.FilePathName = b.String()
}
