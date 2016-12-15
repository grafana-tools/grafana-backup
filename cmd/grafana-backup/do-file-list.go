// Backup tool for Grafana.
// Copyright (C) 2016  Alexander I.Grafov <siberian@laika.name>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// ॐ तारे तुत्तारे तुरे स्व

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grafana-tools/sdk"
)

func doFileList(opts ...option) {
	var (
		cmd  = initCommand(opts...)
		file *os.File
		db   sdk.Board
		ds   sdk.Datasource
		err  error
	)
	for _, filename := range cmd.filenames {
		select {
		case <-cancel:
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
		default:
			if strings.HasSuffix(filename, ".json") {
				if file, err = os.Open(filename); err != nil {
					if cmd.verbose {
						fmt.Fprintf(os.Stderr, "file read error %s\n", err)
					}
					continue
				}
				s := bufio.NewScanner(file)
				s.Split(scanJSONLines)
				for s.Scan() {
					if err = json.Unmarshal(s.Bytes(), &ds); err == nil {
						if ds.Name != "" && ds.URL != "" {
							fmt.Printf("%s:\t source id:%d \"%s\" %s\n", filename, ds.ID, ds.Name, ds.URL)
							continue
						}
					}
					if err = json.Unmarshal(s.Bytes(), &db); err == nil {
						fmt.Printf("%s:\t board id:%d \"%s\"", filename, db.ID, db.Title)
						if len(db.Tags) > 0 {
							fmt.Printf(" %v", db.Tags)
						}
						fmt.Println()
					}
				}
			}
			// if cmd.verbose {
			// 	fmt.Printf("Found %d dashboards.\n", len(foundBoards))
			// }
		}
	}
Exit:
	fmt.Println()
}

func scanJSONLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("}{")); i >= 0 {
		// We have probably a full JSON object followed by another object
		// that came obviously from another file.
		return i + 1, data[0 : i+1], nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
