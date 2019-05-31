package gcf

import (
	"io/ioutil"
	"net/http"

	"github.com/soloworks/go-netlinx/compilelog"
)

// GenerateNetlinxCompileCfg is a Cloud Function which returns a .cfg file for
// Netlinx compiler from a .apw xml file (passed as body)
func GenerateNetlinxCompileCfg(w http.ResponseWriter, r *http.Request) {
	// Get Body as Bytes Array
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Load this into an APW Workspace

	// Get URL Variables
	root := r.URL.Query().Get("root")
	// Process and generate the .cfg
	b, err := compilelog.Process(body, root)

	w.Write(b)
}
