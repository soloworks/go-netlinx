package apw

import "encoding/xml"

// IRDB represetents an AMX project in an APW
type IRDB struct {
	XMLName        xml.Name `xml:"IRDB"`
	Property       string   `xml:"Property"`
	DOSName        string   `xml:"DOSName"`
	UserDBPathName string   `xml:"UserDBPathName"`
	Notes          string   `xml:"Notes"`
	DBKey          string   `xml:"DBKey,attr"`
}
