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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func doRestore(opts ...option) {
	var (
		cmd      = initCommand(opts...)
		rawBoard []byte
		err      error
	)
	for _, filename := range cmd.filenames {
		if strings.HasSuffix(filename, ".json") {
			if rawBoard, err = ioutil.ReadFile(filename); err != nil {
				fmt.Fprintf(os.Stderr, "error on read %s", filename)
				continue
			}
			// TODO add db match filters
			if err = cmd.grafana.SetRawDashboard(rawBoard); err != nil {
				fmt.Fprintf(os.Stderr, "error on importing dashboard from %s", filename)
				continue
			}
			if cmd.verbose {
				fmt.Printf("Dashboard restored from %s.\n", filename)
			}
		}
	}

}
