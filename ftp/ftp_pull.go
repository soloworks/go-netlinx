package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/jlaffaye/ftp" 
)

type amxSys struct {
	Host     string `json:"Host"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type configFile struct {
	DestPath string
	Filename string
	Mask     string   `json:"Mask"`
	AmxSys   []amxSys `json:"Systems"`
}

func main() {

	// Set ConfigFile Variable
	cf := configFile{}

	// Get Command Line Variables
	flag.StringVar(&cf.Filename, "Config", "ftp_pull.json", "Config File name")
	flag.StringVar(&cf.DestPath, "Dest", "files", "Destination File Path")
	flag.Parse()

	// Load in Config Settings
	file, err := os.Open(cf.Filename)
	if err != nil {
		log.Println("Error Loading config file " + cf.Filename)
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

	// Create Waitgroup for sync control
	var wg sync.WaitGroup

	for _, sys := range cf.AmxSys {
		wg.Add(1)
		go getFromAMX(sys.Host+`:21`, sys.Username, sys.Password, cf.Mask, cf.DestPath, &wg)
	}
	wg.Wait()
	log.Println(`Finished!`)
}

func getFromAMX(host string, un string, pw string, mask string, dest string, wg *sync.WaitGroup) error {

	// Dial out to the FTP Server
	log.Println(host + ` Connecting`)
	sc, err := ftp.Dial(host)
	if err != nil {
		wg.Done()
		return err
	}
	// Login with credentials
	log.Println(host + ` Authenticating`)
	err = sc.Login(un, pw)
	if err != nil {
		log.Println(host + ` Login: ` + err.Error())
		wg.Done()
		return err
	}
	// List the files and output to console
	log.Println(host + ` Connected`)
	fs, err := sc.List(``)
	if err != nil {
		wg.Done()
		return err
	}
	for _, f := range fs {
		m, err := filepath.Match(mask, f.Name)
		if err != nil {
			wg.Done()
			return err
		}
		if m {
			log.Println(f.Name)
			r, err := sc.Retr(f.Name)
			if err != nil {
				wg.Done()
				return err
			}
			os.MkdirAll(dest, os.ModePerm)
			outFile, err := os.Create(filepath.Join(dest, f.Name))
			if err != nil {
				wg.Done()
				return err
			}
			log.Println(`Writing: ` + outFile.Name())
			_, err = io.Copy(outFile, r)
			if err != nil {
				wg.Done()
				return err
			}
			outFile.Close()
			r.Close()
		}
	}

	wg.Done()
	return nil
}
