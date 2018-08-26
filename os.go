package apwfile

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"bitbucket.org/solo_works/apwfile/workspace"
)

// FindInFolder searches all subdirectories (recursivly option) for any
// .apw files and returns a list of AMXProjects
func FindInFolder(dir string, recursive bool) []*APWInfo {

	// Create array to hold discovered projects
	var APWInfos []*APWInfo

	// Check root folder and store APWs
	APWInfos = append(APWInfos, getAPWs(dir)...)

	// Do Sub Folder(s) if requested
	if recursive {
		files, _ := ioutil.ReadDir(dir)

		for _, subDir := range files {
			if subDir.IsDir() {
				APWInfos = append(APWInfos, FindInFolder(filepath.Join(dir, subDir.Name()), recursive)...)
			}
		}
	}

	// Return List
	return APWInfos
}
func getAPWs(dir string) []*APWInfo {

	// Create array to hold discovered projects
	var APWInfos []*APWInfo

	// Get all files in possible project folder
	files, _ := ioutil.ReadDir(dir)

	// Cycle through all and identify those with .apw files
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".apw" {
			APWInfos = append(APWInfos, NewAPWInfo(filepath.Join(dir, file.Name())))
		}
	}

	// Return List
	return APWInfos
}

// LoadWorkspace fetches the XML for this apw file
func (apw *APWInfo) LoadWorkspace() error {
	var err error
	apw.Workspace, err = workspace.Load(apw.Filename)
	return err
}

// SaveWorkspace saves the XML to the designated destination folder
func (apw *APWInfo) SaveWorkspace() error {
	workspace.Save(apw.Workspace, apw.Filename)
	return nil
}

// SaveArchive pulls all .apw files together into a zip in the target folder using the workspace name
func (apw *APWInfo) SaveArchive(destDir string, buildID string) error {
	// Verify the APW file is all good before we do this
	missing, err := apw.checkFiles()
	if err != nil {
		var e bytes.Buffer
		e.WriteString(strconv.Itoa(len(missing)))
		e.WriteString(" File")
		if len(missing) > 0 {
			e.WriteString("s")
		}
		e.WriteString(" not found")
		return errors.New(e.String())
	}

	// Create a Zip writer with defered close
	var filename bytes.Buffer
	filename.WriteString(apw.Identifier)
	if buildID != "" {
		filename.WriteString("_" + buildID)
	}
	filename.WriteString(".zip")
	myZipFile, err := os.Create(filepath.Join(destDir, filename.String()))
	if err != nil {
		return err
	}
	defer myZipFile.Close()

	z := zip.NewWriter(myZipFile)
	defer z.Close()

	// Add each file to the Archive
	for file, fileType := range apw.Files {

		// Open existing file
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)

		// Set file to correct folder based on file type
		header.Name = filepath.Join(workspace.FileFolder(fileType), header.Name)

		// Compress File
		header.Method = zip.Deflate

		writer, err := z.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}

	}

	// Save XML
	f, err := z.Create(apw.Identifier + ".apw")
	if err != nil {
		log.Fatal(err)
	}
	myXML, _ := apw.Workspace.Bytes()
	_, err = f.Write([]byte(myXML))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
