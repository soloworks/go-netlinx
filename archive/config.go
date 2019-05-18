package main

import (
	"encoding/json"
	"os"
	"regexp"
)

// Configuration holds imported json config file data
type Configuration struct {
	CompilerPath    string `json:"compilerPath"`
	RepoRootPath    string `json:"repoRootPath"`
	GlobalAMXFolder string `json:"globalAMXFolder"`
	GlobalProject   string `json:"globalProject"` 
}

// Load imports from a JSON file
func Load(configFile string) Configuration {
	// Load in Config Settings
	file, fileErr := os.Open(configFile)
	if fileErr != nil {
		os.Exit(0)
	}

	// Decode Config Settings
	decoder := json.NewDecoder(file)
	config := Configuration{}
	decodeErr := decoder.Decode(&config)
	if decodeErr != nil {
		os.Exit(1)
	}

	// Replace any environment variables referenced in the config file
	// Set a new reg to look for userprofile tag in case insensitive mode
	re := regexp.MustCompile(`(?i)%userprofile%`)

	// Replace all instances with the actual environment variable
	config.RepoRootPath = re.ReplaceAllString(config.RepoRootPath, os.Getenv("USERPROFILE"))

	return config
}
