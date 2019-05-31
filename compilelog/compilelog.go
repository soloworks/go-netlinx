package compilelog

import (
	"bufio"
	"bytes"
	"fmt"
)

// ProcessLog takes a raw log file generated via Netlinx Compiler
// And returns a more readable version with stats
func ProcessLog(logData []byte) ([]byte, error) {

	// Put the bytes into a scanner
	scanner := bufio.NewScanner(bytes.NewReader(logData))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil, nil
}
