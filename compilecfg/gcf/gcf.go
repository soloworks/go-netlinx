package gcf

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/soloworks/go-netlinx/apw"
	"github.com/soloworks/go-netlinx/compilecfg"
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
	a, err := apw.NewAPW("myWorkspace.apw", body)

	// Get URL Variables
	root := r.URL.Query().Get("root")
	logfile := r.URL.Query().Get("logfile")
	logconsole, _ := strconv.ParseBool(r.URL.Query().Get("logconsole"))
	// Process and generate the .cfg
	b := compile.GenerateCFG(*a, root, logfile, logconsole)

	w.Write(b)
}
