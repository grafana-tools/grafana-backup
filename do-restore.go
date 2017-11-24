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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/grafana-tools/sdk"
)

//THIS IS BROKEN
//It no longer restores
//
//The last thing I did was change from CreateDatasource to UpdateDatasource. Perhaps I did not recompile after I made the
//change and I tested the wrong code (which is why I need to start using 'go run *.go' instead of compiling at all.)
//
//Change it back.
//
//Then figure out what happens if the datasource exists and I use CreateDatasource.
//
//Also it appears that restoreDatasources is being run twice.

// Triggers a restore.
func doRestore(opts ...option) {
	var (
		cmd      = initCommand(opts...)
	)

	// TODO: apply-to auto doesn't make much sense in the context of restore yet.
	// An actual heirarchal restore probably isn't feasable until we get better about parsing the JSON.
	// There isn't going to be a guarantee that the datasource filename is exactly what we expect.
	if cmd.applyHierarchically {
		restoreDashboards(cmd)
		restoreDatasources(cmd)
		restoreUsers(cmd)
		return
	}
	if cmd.applyForBoards {
		restoreDashboards(cmd)
	}
	if cmd.applyForDs {
		//restoreDatasources(cmd, nil)
		restoreDatasources(cmd)
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
		//datasources = make(map[string]bool) // If cmd.applyHierarchically is true extract datasources from the dashboard and restore those as well.
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
		} //else {
		//	if cmd.verbose {
		//		fmt.Fprintf(os.Stderr, "File %s does not appear to be a dashboard: Skipping file.", filename)
		//	}
		//
		//}
	}

	// Disabling the 'heirarchal' functionality until it can be implemented properly.
	//if cmd.applyHierarchically {
	//	restoreDatasources(cmd)
	//}
}

// Restores all datasource files. Currently those are files that match the format .*.ds.([0-9]+).json.
func restoreDatasources(cmd *command) {
	var (
		rawDS          []byte
		err            error
	)

	for _, filename := range cmd.filenames {
		pattern, _ := regexp.Compile(".*.ds.([0-9]+).json")

		if pattern.MatchString(filename) {
			if rawDS, err = ioutil.ReadFile(filename); err != nil {
				fmt.Fprintf(os.Stderr, "error on read %s", filename)
				continue
			}

			// TODO: most of this should probably be pushed upstream into grafana SDK in a CreateRawDatasource function
			// Stolen from SetRawDashboard
			var (
				resp    sdk.StatusMessage
				err     error
				plain   sdk.Datasource
			)

			if err = json.Unmarshal(rawDS, &plain); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling datasource from file %s: %s\n", filename, err)
				continue
			}

			// TODO: Check to see if the datasource already exists and use the correct method or throw an error on update unless --force is specified.
			resp, err = cmd.grafana.CreateDatasource(plain)
			//resp, err = cmd.grafana.UpdateDatasource(plain)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error importing datasource from %s: %s\n", filename, err)
				continue
			}

			if *resp.Message == "Data source with same name already exists" {
				//TODO: Update this so that we pull out the datasource name and give that in the message.
				fmt.Fprintf(os.Stderr, "A Datasource with the same name as specified in %s already exists.\n", filename)
				continue
			} else if *resp.Message != "Datasource added" {
				fmt.Fprintf(os.Stderr, "Error importing datasource from %s: %s\n", filename, *resp.Message)
				continue
			}

			if cmd.verbose {
				fmt.Printf("Datasource restored from %s.\n", filename)
			}
		} //else {
		//	if cmd.verbose {
		//		fmt.Fprintf(os.Stderr, "File %s does not appear to be a datasource: Skipping file.\n", filename)
		//	}
		//
		//}
	}
}

// Not yet implemented.
func restoreUsers(cmd *command) {
	if cmd.verbose {
		fmt.Fprintln(os.Stderr, "Restoring users not yet implemented!")
	}
}
