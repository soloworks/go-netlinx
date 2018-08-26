package devicemap

import "encoding/xml"

// DeviceMap represetents an AMX project in an APW
type DeviceMap struct {
	XMLName xml.Name `xml:"DeviceMap"`
	DevAddr string   `xml:"DevAddr,attr"`
	DevName string   `xml:"DevName,omitempty"`
}

// NewDeviceMap returns a new project instance with
// default fields already populated
func NewDeviceMap(devAddr string, devName string) *DeviceMap {
	return &DeviceMap{
		DevName: devName,
		DevAddr: devAddr,
	}
}
