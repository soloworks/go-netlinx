package compilelog

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

// Process takes a raw log file generated via Netlinx Compiler
// And returns a more readable version with stats
func Process(logData []byte, root string) ([]byte, error) {

	// Variables to hold data
	var FilesTotal int
	var FilesTotalWarning int
	var FilesTotalError int
	var FileWarning bool
	var FileError bool
	var FileName string
	// String Builder
	var sb strings.Builder

	// Put the bytes into a scanner
	scanner := bufio.NewScanner(bytes.NewReader(logData))
	for scanner.Scan() {
		// Get a line
		l := scanner.Text()
		// Trim Whitespace
		l = strings.TrimSpace(l)
		// Remove root if present
		l = strings.Replace(l, root, "", 1)
		// Only handle lines with more than 1 char
		if len(l) > 0 {
			if strings.HasPrefix(l, "---- Starting NetLinx Compile") {
				FileWarning = false
				FileError = false
				FileName = ""
			}
			if strings.HasPrefix(l, "WARNING:") || strings.HasPrefix(l, "ERROR:") {
				if strings.HasPrefix(l, "WARNING:") {
					FileWarning = true
					l = strings.TrimPrefix(l, "WARNING:")

				}
				if strings.HasPrefix(l, "ERROR:") {
					FileError = true
					l = strings.TrimPrefix(l, "ERROR:")
				}
				l = strings.TrimSpace(l)

				chunks := strings.SplitN(l, "(", 2)
				if FileName == "" {
					FileName = chunks[0]
					sb.WriteString("\n")
					sb.WriteString(FileName)
					sb.WriteString("\n")
				}

				sb.WriteString("(")
				sb.WriteString(chunks[1])
				sb.WriteString("\n")

			}
			if strings.HasPrefix(l, "---- NetLinx Compile Complete") {
				if FileWarning {
					FilesTotalWarning++
				}
				if FileError {
					FilesTotalError++
				}
				FilesTotal++
			}
		}
	}

	// Add the Stats
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString("Files with Warnings: ")
	sb.WriteString(strconv.Itoa(FilesTotalWarning))
	sb.WriteString("\n")
	sb.WriteString("Files with Errors: ")
	sb.WriteString(strconv.Itoa(FilesTotalError))
	sb.WriteString("\n")
	sb.WriteString("Files Processed: ")
	sb.WriteString(strconv.Itoa(FilesTotal))
	sb.WriteString("\n")

	return []byte(sb.String()), nil
}
