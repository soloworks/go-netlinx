package system

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/soloworks/apwfile/file"
	"github.com/soloworks/apwfile/transport"
)

// System represetents an AMX project in an APW
type System struct {
	XMLName                  xml.Name     `xml:"System"`
	Identifier               string       `xml:"Identifier"`
	SysID                    int          `xml:"SysID"`
	TransTCPIP               string       `xml:"TransTCPIP,omitempty"`
	TransSerial              string       `xml:"TransSerial,omitempty"`
	TransTCPIPEx             string       `xml:"TransTCPIPEx,omitempty"`
	TransSerialEx            string       `xml:"TransSerialEx,omitempty"`
	TransUSBEx               string       `xml:"TransUSBEx,omitempty"`
	TransVNMEx               string       `xml:"TransVNMEx,omitempty"`
	VirtualNetLinxMasterFlag string       `xml:"VirtualNetLinxMasterFlag,omitempty"`
	VNMSystemID              string       `xml:"VNMSystemID,omitempty"`
	VNMIPAddress             string       `xml:"VNMIPAddress,omitempty"`
	VNMMaskAddress           string       `xml:"VNMMaskAddress,omitempty"`
	UserName                 string       `xml:"UserName,omitempty"`
	Password                 string       `xml:"Password,omitempty"`
	Comments                 string       `xml:"Comments,omitempty"`
	Files                    []*file.File `xml:"File"`
	IsActive                 string       `xml:"IsActive,attr"`
	Platform                 string       `xml:"Platform,attr"`
	Transport                string       `xml:"Transport,attr"`
	TransportEx              string       `xml:"TransportEx,attr"`
}

// ByIdentifier implements sort.Interface for []Project based on
// the Identifier field.
type ByIdentifier []*System

func (a ByIdentifier) Len() int           { return len(a) }
func (a ByIdentifier) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIdentifier) Less(i, j int) bool { return a[i].Identifier < a[j].Identifier }

// NewSystem returns a new project instance with
// default fields already populated
func NewSystem(identifier string, sysID int) *System {
	var s System
	s.SysID = sysID
	s.Identifier = fmt.Sprintf("%03d", sysID)
	s.Identifier += ": " + identifier
	return &s
}

// AddConnectionToSystem adds and sets an IP connection to the system
func (s *System) AddConnectionToSystem(t *transport.Transport) {
	// Set the Type
	s.TransportEx = t.Type

	// Create new buffer
	var buf bytes.Buffer
	// Concat each value with demilter of |
	buf.WriteString(t.Host + "|")
	buf.WriteString(strconv.Itoa(t.Port) + "|")
	// Convert Bool
	if t.PingTest {
		buf.WriteString("1|")
	} else {
		buf.WriteString("0|")
	}

	buf.WriteString(t.Name + "|")
	buf.WriteString(t.Username + "|")
	buf.WriteString(t.Password)
	// Store the value
	s.TransTCPIPEx = buf.String()
}

// FindFile returns a pointer to a system
func (s *System) FindFile(id string) *file.File {
	for i, f := range s.Files {
		if f.Identifier == id {
			return s.Files[i]
		}
	}
	return nil
}

// AddFile adds a file to a system
func (s *System) AddFile(f *file.File) {
	s.Files = append(s.Files, f)
}
