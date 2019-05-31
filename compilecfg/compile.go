package compilecfg

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/soloworks/go-netlinx/apw"
)

func toWindows(x string) string {
	return strings.Join(strings.Split(x, `/`), `\`)
}

func toLinux(x string) string {
	return strings.Join(strings.Split(x, `\`), `/`)
}

// GenerateCFG creates a Netlinx Compiler .cfg file from a workspace
func GenerateCFG(a apw.APW, root string, logfile string, logconsole bool) []byte {

	// Create an empty list for modules
	var Modules []string
	var Source []string
	IncludePath := make(map[string]struct{})

	// Extract list of .axs Modules and .axs Source
	for x, y := range a.FilesReferenced {
		switch filepath.Ext(x) {
		case ".axs", ".axi":
			switch y {
			case "Module":
				Modules = append(Modules, x)
			case "Source", "MasterSrc":
				Source = append(Source, x)
			case "Include":
				//root = strings.ReplaceAll(root, `\`, `/`)
				//root = toLinux(root)
				//x = strings.ReplaceAll(x, `\`, `/`)
				x = toLinux(x)
				ip := filepath.Dir(x)
				//ip = strings.ReplaceAll(ip, `/`, `\`)
				ip = toWindows(ip)
				IncludePath[ip] = struct{}{}
			}
		}
	}

	// Order the lists
	sort.Strings(Modules)
	sort.Strings(Source)

	// Fix an Unix folders to Windows
	root = toWindows(root)

	// Build the Config File Header & Options
	var sb strings.Builder
	sb.WriteString(";------------------------------------------------------------------------------\n")
	sb.WriteString(";\n")
	sb.WriteString("; Netlinx Compiler Config File generated by Go\n")
	sb.WriteString("; Source: http://github.org/soloworks/go-netlinx/compiler\n")
	sb.WriteString(`; Run> NLRC -C"`)
	if root != "" {
		sb.WriteString(root)
		sb.WriteString(`\`)
	}
	sb.WriteString(`filename.cfg"`)
	sb.WriteString(";\n")
	sb.WriteString(";------------------------------------------------------------------------------\n\n")

	// Write out Root Directory
	if root == "" {
		sb.WriteString("MainAXSRootDirectory=-R\n\n")
	} else {
		sb.WriteString("MainAXSRootDirectory=")
		sb.WriteString(root)
		sb.WriteString("\n\n")
	}
	if logfile != "" {
		sb.WriteString("OutputLogFileOption=N\n")
		sb.WriteString("OutputLogFile=")
		sb.WriteString(root)
		sb.WriteString("\\")
		sb.WriteString(logfile)
		sb.WriteString("\n")
	}
	if logconsole {
		sb.WriteString("OutputLogConsoleOption=Y\n")
	} else {
		sb.WriteString("OutputLogConsoleOption=N\n")
	}
	sb.WriteString("BuildWithDebugInformation=Y\n")
	sb.WriteString("BuildWithSource=N\n")
	sb.WriteString("BuildWithWC=Y\n\n")

	// Add the Include File Folders
	for x := range IncludePath {
		sb.WriteString("AdditionalIncludePath=")
		sb.WriteString(x)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	// Add the Modules
	for _, x := range Modules {
		sb.WriteString("AXSFILE=")
		sb.WriteString(x)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// Add the Source
	for _, x := range Source {
		sb.WriteString("AXSFILE=")
		sb.WriteString(x)
		sb.WriteString("\n")
	}
	// Return the file as bytes
	return []byte(sb.String())
}
