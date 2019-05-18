# netlinxstudioconfigurator : Go based tool for configuring AMX Netlinx Studio installations

## Overview [![GoDoc](https://godoc.org/bitbucket.org/solo_works/netlinxstudioconfigurator?status.svg)](https://godoc.org/bitbucket.org/solo_works/netlinxstudioconfigurator)

This package produces a tool to allow automatic configuration of various settings for AMX Netlinx Studio. These include setting of default folders for Module, Include and Library files for centralised code.

It was created to automate the process of pointing new installations at our core files.

## Compile and Run

Clone the repo and build using Go Build

## Download PreCompiled

Versioned Binaries are in the near future. For now there is a copy of the latest binary in the Downloads section of this repo

## Command Line Settings

Clear all existing Folder settings (Keeps all defaults i.e. with AMXShare in the path)
```console
> netlinxstudioconfigurator -Clear 
```

Specify a config file for use (Default: config.json)
```console
> netlinxstudioconfigurator -Config xxx.json 
```

## Configuration File
Sample config file (config.json)
```json
{
  "BasePath": "%USERPROFILE%\\Documents\\GitHub\\netlinx-global-code",
  "Modules": [
    "ModulesDuet",
    "ModulesNetlinx",
    "ModulesDuetRms",
    "ModulesNetlinxRms"
  ],
  "Includes": [
    "Includes",
    "IncludesRms"
  ],
  "Libs": [],
  "CompileWithDebug": true,
  "CompileWithSrc": false,
  "SmartTransfer": true,
  "SendSource": false,
  "TabWidth": 3,
  "IndentWidth": 3
}
```

## Windows Registry

The application will update registry keys in thej following locations:

*  Computer\\HKEY_LOCAL_MACHINE\\SOFTWARE\\WOW6432Node\\AMX Corp.\\NetLinx Studio\\NLXCompiler_Includes
*  Computer\\HKEY_LOCAL_MACHINE\\SOFTWARE\\WOW6432Node\\AMX Corp.\\NetLinx Studio\\NLXCompiler_Modules
*  Computer\\HKEY_CURRENT_USER\\Software\\AMX Corp.\\NetLinx Studio\\Editor Preferences
*  Computer\\HKEY_CURRENT_USER\\Software\\AMX Corp.\\NetLinx Studio\\NLXCompiler_Options
*  Computer\\HKEY_CURRENT_USER\\Software\\AMX Corp.\\NetLinx Studio\\Batch Transfer User Options

## Author

Created by Sam Shelton for Solo Works London

Find us at https://soloworks.co.uk/