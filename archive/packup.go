package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time" 
	"github.com/soloworks/go-netlinx/apw"
)

type myargs struct {
	path      string
	workspace string 
	project   string
	system    string
	archive   bool
	handover  bool
	release   bool
}

func main() {

	// Local Variables
	c := Load("packup.json") 
	var args myargs

	// Specify Build ID from UTC
	buildID := time.Now().UTC().Format("2006_01_02_")
	buildID = buildID + strconv.Itoa(time.Now().Hour()*3600+time.Now().Minute()*60+time.Now().Second())
	//buildID := time.Now().UTC().Format("20060102150405")
	//buildID := strconv.Itoa(int(time.Now().Unix()))

	// Variables used to build project
	var w = -1
	var p = -2
	var s = -2

	// Get Command Line Variables
	flag.StringVar(&args.path, "Working Path", "", "Path to use (overrides default in config)")
	flag.StringVar(&args.workspace, "Workspace", "", "Use the specified APW file if found")
	flag.StringVar(&args.project, "Project", "All", "Specified project to process (Default = All)")
	flag.StringVar(&args.system, "System", "All", "Specified System to process (Default = All)")
	flag.BoolVar(&args.archive, "A", false, "Produce Archive Package")
	flag.BoolVar(&args.handover, "H", false, "Produce Handover Package")
	flag.BoolVar(&args.release, "R", false, "Produce Release Package")

	flag.Parse()

	// Process any command line stuff
	if args.path != "" {
		c.RepoRootPath = args.path
	}
	if args.project == "All" {
		p = -1
		s = -1
	}
	if args.system == "All" {
		s = -1
	}
	if !args.archive && !args.handover && !args.release {
		fmt.Println("No package type set - use -A|H|R")
		os.Exit(3)
	}

	// Debug the argumentgs
	fmt.Println("Command Line args ", args)

	// Get AMX Projects
	apwInfos := apwfile.FindInDir(c.RepoRootPath)

	// Check for command line workspace match
	for index, apwInfo := range apwInfos {
		if apwInfo.Identifier == args.workspace {
			// Set this as the working workspace
			w = index
			// Exit Loop
			break
		}
	}

	// Output workspaces for user input
	if w == -1 {
		for index, apwInfo := range apwInfos {
			fmt.Printf("%02d", index)
			fmt.Println(": ", apwInfo.Identifier)
		}
		// Prompt User
		w = cli.AskForItemIndex("Select Workspace:", false)
		fmt.Println("Selected: ", apwInfos[w].Name)
	}

	// Load up the requested workspace
	apwfile.LoadWorkspace(&apwInfos[w])

	// Check for command line project match
	if args.project != "" {
		for index, project := range apwInfos[w].XML.Projects {
			if project.Identifier == args.project {
				// Set this as the working workspace
				p = index
				// Exit Loop
				break
			}
		}
	}
	// Check for command line system match
	if p > -1 && args.system != "" {
		for index, system := range apwInfos[w].XML.Projects[p].Systems {
			if system.Identifier == args.system {
				// Set this as the working workspace
				s = index
				// Exit Loop
				break
			}
		}
	}

	// Prompt for a project if not set
	if p == -2 {
		for index, project := range apwInfos[w].XML.Projects {
			fmt.Printf("%02d", index)
			fmt.Println(": " + project.Identifier)
		}
		p = cli.AskForItemIndex("Select Project (or a for All):", true)
	}

	// Process project
	if p > -1 {
		// Swap array for new array with a single project
		apwInfos[w].XML.Projects = []workspace.Project{apwInfos[w].XML.Projects[p]}
	}

	// If one project has been selected, prompt for system
	if p > -1 && s == -2 {

		for index, system := range apwInfos[w].XML.Projects[p].Systems {
			fmt.Printf("%02d", index)
			fmt.Println(": " + system.Identifier)
		}
		s = cli.AskForItemIndex("Select System (or a for All):", true)
		if s > -1 {
			apwInfos[w].XML.Projects[p].Systems = []workspace.System{apwInfos[w].XML.Projects[p].Systems[s]}
		}
	}

	// Process an Archive Package
	if args.archive {
		packItUp(apwInfos[w], ptArchive, buildID)
	}

	// Process a Handover Package
	if args.handover {
		packItUp(apwInfos[w], ptHandover, buildID)
	}

	// Process a Release Package
	if args.release {
		packItUp(apwInfos[w], ptRelease, buildID)
	}

	log.Println("Done!")
}

type packType int

const (
	ptArchive packType = iota
	ptHandover
	ptRelease
)

func packItUp(apw apwfile.APWInfo, pt packType, buildID string) {

	// Ammend the BuildID with the type of packaging being done
	switch pt {
	case ptArchive:
		log.Println("Packing Archive")
		buildID = buildID + "A"
	case ptHandover:
		log.Println("Packing Handover")
		buildID = buildID + "H"
	case ptRelease:
		log.Println("Packing Release")
		buildID = buildID + "R"
	}

	//log.Println(apw.XML)

	// ADjust the XML to reflect the package type
	switch pt {
	case ptHandover:
		workspace.ConvertToHandover(&apw.XML)
	case ptRelease:
		workspace.ConvertToRelease(&apw.XML)
	}

	// Populate the Files referece
	apwfile.GetFiles(&apw)

	// Adjust the XML to relative paths
	workspace.SetRelativePaths(&apw.XML)

	// get the full path of a release folder
	destDir := filepath.Join(apw.OriginPath, "PackUp")
	// Create it ignoring any errors (such as it existing already)
	_ = os.Mkdir(destDir, os.ModePerm)

	// Zip it up!
	err := apwfile.Zip(&apw, destDir, buildID)
	if err != nil {

	}

	// Announce Finish
	switch pt {
	case ptArchive:
		log.Println("Completed Archive")
	case ptHandover:
		log.Println("Completed Handover")
	case ptRelease:
		log.Println("Completed Release")
	}
}
