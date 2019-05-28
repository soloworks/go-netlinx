package cloud

import (
	"io/ioutil"
	"net/http"

	"github.com/soloworks/go-netlinx/apw"
	"github.com/soloworks/go-netlinx/compile"
)

func cloudFunction(w http.ResponseWriter, r *http.Request) {
	// Get Body as Bytes Array
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Load this into an APW Workspace
	a, err := apw.NewAPW("myWorkspace.apw", body)

	// Process and generate the .cfg
	b := compile.GenerateCFG(*a)

	w.Write(b)
}
