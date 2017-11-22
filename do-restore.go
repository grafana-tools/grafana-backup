// Backup tool for Grafana.
// Copyright (C) 2016-2017 Alexander I.Grafov <siberian@laika.name>
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
		//rawBoard []byte
		//err      error
	)

	// A heirarchal restore probably isn't feasable until we get better about parsing the JSON. There isn't
	// going to be a guarantee that the datasource filename is exactly what we expect.
	if cmd.applyHierarchically {
		restoreDashboards(cmd)
		return
	}
	if cmd.applyForBoards {
		restoreDashboards(cmd)
	}
	if cmd.applyForDs {
		restoreDatasources(cmd, nil)
	}
	if cmd.applyForUsers {
		restoreUsers(cmd)
	}

}

// Restores all dashboard files. Currently that's files that end in .db.json
// Then if cmd.applyHierarchically is true calls restoreDatasources
func restoreDashboards(cmd *command) {
	var (
		rawBoard    []byte
		datasources = make(map[string]bool) // If cmd.applyHierarchically is true extract datasources from the dashboard and restore those as well.
		err         error
		// These three are used in backupDashboards, figure out what they're used for and if I want to implement them. -AF
		//boardLinks  []sdk.FoundBoard
		//meta        sdk.BoardProperties
		//board       sdk.Board
	)

	for _, filename := range cmd.filenames {
		if strings.HasSuffix(filename, "db.json") {
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
		} else {
			if cmd.verbose {
				fmt.Fprintf(os.Stderr, "File %s does not appear to be a dashboard: Skipping file.", filename)
			}

		}
	}


	if cmd.applyHierarchically {
		restoreDatasources(cmd, datasources)
	}
}

func restoreDatasources(cmd *command, datasources map[string]bool) {
	if cmd.verbose {
		fmt.Fprint(os.Stderr, "Restoring datasources not yet implemented!")
	}
}

func restoreUsers(cmd *command) {
	if cmd.verbose {
		fmt.Fprint(os.Stderr, "Restoring users not yet implemented!")
	}
}




// Extracts a map of datasources used by a dashboard. This function currently exists inside of do-backup
// and is here because I have to figure out how to use it when doing a restore then if I use the exact
// function probably move it to a common library.
//func extractDatasources(datasources map[string]bool, board sdk.Board) {
//	for _, row := range board.Rows {
//		for _, panel := range row.Panels {
//			if panel.Datasource != nil {
//				datasources[*panel.Datasource] = true
//				fmt.Println(slug.Make(*panel.Datasource))
//			}
//		}
//	}
//}