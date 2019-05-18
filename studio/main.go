package main 

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sys/windows/registry" 
)

// ConfigFile is a structure holding all local settings
type ConfigFile struct {
	CompileDebug  bool     `json:"CompileWithDebug"` 
	CompileSrc    bool     `json:"CompileWithSrc"`
	SendSource    bool     `json:"SendSource"`
	SmartTransfer bool     `json:"SmartTransfer"`
	TabWidth      uint32   `json:"TabWidth"`
	IndentWidth   uint32   `json:"IndentWidth"`
	BasePath      string   `json:"BasePath"`
	Modules       []string `json:"Modules"`
	Includes      []string `json:"Includes"`
	Libs          []string `json:"Libs"`
}

type regEntry struct {
	Name  string
	Value interface{}
}

var version = "undefined"

type myargs struct {
	configFile    string
	clearExisting bool
}

var args myargs

func main() {
	// Set ConfigFile Variable
	cf := ConfigFile{}

	// Get Command Line Variables
	flag.StringVar(&args.configFile, "Config", "config.json", "Config File name")
	flag.BoolVar(&args.clearExisting, "Clear", false, "Clear Existing Settings")
	flag.Parse()

	if !args.clearExisting {
		// Load in Config Settings
		file, err := os.Open(args.configFile)
		if err != nil {
			log.Println("Error Loading config file " + args.configFile)
			log.Println(err)
			os.Exit(0)
		}

		// Read Config Settings
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&cf)
		if err != nil {
			log.Println("Config File Error: ")
			log.Println(err)
			os.Exit(0)
		}

		// Alter Config settings for %USERPROFILE% if present
		cf.BasePath = strings.Replace(cf.BasePath, `%USERPROFILE%`, os.Getenv("USERPROFILE"), 1)

		// Change file paths to fully qualified based on base folder
		for i, stub := range cf.Includes {
			cf.Includes[i] = filepath.Join(cf.BasePath, stub)
		}
		for i, stub := range cf.Modules {
			cf.Modules[i] = filepath.Join(cf.BasePath, stub)
		}
		for i, stub := range cf.Libs {
			cf.Libs[i] = filepath.Join(cf.BasePath, stub)
		}
	}
	// Create Waitgroup for sync control
	var wg sync.WaitGroup

	// Call Folder Function
	wg.Add(3)
	updateFolders("Include", cf.Includes, &wg)
	updateFolders("Module", cf.Modules, &wg)
	updateFolders("Lib", cf.Libs, &wg)

	if !args.clearExisting {
		// Update Settings - Editor Preferences
		Values := []regEntry{}
		if cf.IndentWidth > 0 {
			Values = append(Values, regEntry{"IndentWidth", cf.IndentWidth})
		}
		if cf.TabWidth > 0 {
			Values = append(Values, regEntry{"TabWidth", cf.TabWidth})
		}
		wg.Add(1)
		updateSettings("Editor Preferences", Values, &wg)

		// Update Settings - Batch Transfer User Options
		Values = []regEntry{}
		Values = append(Values, regEntry{"TP4 Smart Transfer", cf.SmartTransfer})
		Values = append(Values, regEntry{"TP5 Smart Transfer", cf.SmartTransfer})
		Values = append(Values, regEntry{"Auto Send SRC", cf.SendSource})
		wg.Add(1)
		updateSettings("Batch Transfer User Options", Values, &wg)

		// Update Settings - NLXCompiler_Options
		Values = []regEntry{}
		Values = append(Values, regEntry{"BuildWithDebugInfo", cf.CompileDebug})
		Values = append(Values, regEntry{"BuildWithSource", cf.CompileSrc})
		Values = append(Values, regEntry{"EnableWC", true})
		wg.Add(1)
		updateSettings("NLXCompiler_Options", Values, &wg)
	}
	// Wait, Tidyup, End
	wg.Wait()
	if args.clearExisting {
		fmt.Println("Configuration Cleared")
	} else {
		fmt.Println("Configuration Applied")
	}
}

func updateSettings(KeyName string, Values []regEntry, wg *sync.WaitGroup) {
	// Defer Sync
	defer fmt.Println(`Updated Registry: HKEY_CURRENT_USER\Software\AMX Corp.\NetLinx Studio\` + KeyName)
	defer wg.Done()

	// Open Include Files Key
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\AMX Corp.\NetLinx Studio\`+KeyName, registry.ALL_ACCESS)
	if err != nil {
		log.Println("Error opening key: ")
		log.Fatal(err)
	}
	defer k.Close()

	// Set all of the Values
	for _, v := range Values {
		var newValue uint32
		switch v.Value.(type) {
		case bool:
			if v.Value == true {
				newValue = 1
			}
		case uint32:
			newValue = v.Value.(uint32)
		}
		// Write Value to ValueName key
		k.SetDWordValue(v.Name, newValue)
	}

}

func updateFolders(FolderType string, Folders []string, wg *sync.WaitGroup) {
	// Defer Sync
	defer fmt.Println(`Updated Registry: HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\NLXCompiler_` + FolderType)
	defer wg.Done()

	// Open Include Files Key
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\AMX Corp.\NetLinx Studio\NLXCompiler_`+FolderType+`s`, registry.ALL_ACCESS)
	if err != nil {
		log.Println("Error opening key: ")
		log.Fatal(err)
	}
	defer k.Close()

	// Read the list of value names
	ValueNames, err := k.ReadValueNames(0)
	if err != nil {
		log.Println("Error Reading ValueNames: ")
		log.Fatal(err)
	}

	// Process Existing Value strings
	var vd []string
	for _, vn := range ValueNames {
		// Get the Value
		v, _, _ := k.GetStringValue(vn)
		// Append or Prepend if standard or user
		if strings.Contains(v, `Common Files\AMXShare`) {
			// Put system defaults to front of the list
			vd = append([]string{v}, vd...)
		} else {
			// Put custom to end of the list (or not if clear enabled)
			if !args.clearExisting {
				vd = append(vd, v)
			}
		}
		// Delete value as we will re-add them in a moment
		err := k.DeleteValue(vn)
		if err != nil {
			log.Println("Error Deleting Key:")
			log.Println(err)

		}
	}

	// Add our new ones (if not present and not set to clear)
	if !args.clearExisting {
		for _, f := range Folders {
			var found bool
			for x := range vd {
				if f == vd[x] {
					found = true
				}
			}
			if !found {
				vd = append(vd, f)
			}
		}
	}

	// Store all keys back to the Registry
	for x := range vd {
		// Build Value Name
		var s bytes.Buffer
		switch FolderType {
		case "Lib":
			s.WriteString("Library")
		default:
			s.WriteString(FolderType)
		}
		s.WriteString("Directory")
		s.WriteString(fmt.Sprintf("%03d", x))
		// Write Value to ValueName key
		k.SetStringValue(s.String(), vd[x])
	}
}
