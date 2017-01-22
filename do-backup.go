// Backup tool for Grafana.
// Copyright (C) 2016-2017  Alexander I.Grafov <siberian@laika.name>
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

	"github.com/gosimple/slug"
	"github.com/grafana-tools/sdk"
)

func doBackup(opts ...option) {
	var (
		cmd = initCommand(opts...)
	)
	if cmd.applyHierarchically {
		backupDashboardsHierchically(cmd)
		return
	}
	if cmd.applyForBoards {
		backupDashboards(cmd)
	}
	if cmd.applyForDs {
		backupDatasources(cmd, nil)
	}
	if cmd.applyForUsers {
		backupUsers(cmd)
	}

}

// TODO merge with backupDashboards
func backupDashboardsHierchically(cmd *command) {
	var (
		boardLinks  []sdk.FoundBoard
		rawBoard    []byte
		meta        sdk.BoardProperties
		board       sdk.Board
		datasources = make(map[string]bool)
		err         error
	)
	if boardLinks, err = cmd.grafana.SearchDashboards(cmd.boardTitle, cmd.starred, cmd.tags...); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	if cmd.verbose {
		fmt.Printf("Found %d dashboards that matched the conditions.\n", len(boardLinks))
	}
	for _, link := range boardLinks {
		select {
		case <-cancel:
			exitBySignal()
		default:
			if rawBoard, meta, err = cmd.grafana.GetRawDashboard(link.URI); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, link.URI))
				continue
			}
			if err = json.Unmarshal(rawBoard, &board); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("error %s parsing %s\n", err, meta.Slug))
			} else {
				extractDatasources(datasources, board)
			}
			var fname = fmt.Sprintf("%s.db.json", meta.Slug)
			if err = ioutil.WriteFile(fname, rawBoard, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, meta.Slug))
				continue
			}
			if cmd.verbose {
				fmt.Printf("%s writen into %s.\n", meta.Slug, fname)
			}
		}
	}
	backupDatasources(cmd, datasources)
}

func backupDashboards(cmd *command) {
	var (
		boardLinks []sdk.FoundBoard
		rawBoard   []byte
		meta       sdk.BoardProperties
		err        error
	)
	if boardLinks, err = cmd.grafana.SearchDashboards(cmd.boardTitle, cmd.starred, cmd.tags...); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	if cmd.verbose {
		fmt.Printf("Found %d dashboards that matched the conditions.\n", len(boardLinks))
	}
	for _, link := range boardLinks {
		select {
		case <-cancel:
			exitBySignal()
		default:
			if rawBoard, meta, err = cmd.grafana.GetRawDashboard(link.URI); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, link.URI))
				continue
			}
			var fname = fmt.Sprintf("%s.db.json", meta.Slug)
			if err = ioutil.WriteFile(fname, rawBoard, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, meta.Slug))
				continue
			}
			if cmd.verbose {
				fmt.Printf("%s writen into %s.\n", meta.Slug, fname)
			}
		}
	}
}

func backupUsers(cmd *command) {
	var (
		allUsers []sdk.User
		rawUser  []byte
		err      error
	)
	if allUsers, err = cmd.grafana.GetAllUsers(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		return
	}
	for _, user := range allUsers {
		select {
		case <-cancel:
			exitBySignal()
		default:
			rawUser, _ = json.Marshal(user)
			var fname = fmt.Sprintf("%s.user.%d.json", slug.Make(user.Login), user.OrgID)
			if err = ioutil.WriteFile(fname, rawUser, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("error %s on writing %s\n", err, fname))
				continue
			}
			if cmd.verbose {
				fmt.Printf("%s written into %s\n", user.Name, fname)
			}
		}
	}
}

func backupDatasources(cmd *command, datasources map[string]bool) {
	var (
		allDatasources []sdk.Datasource
		rawDs          []byte
		err            error
	)
	if allDatasources, err = cmd.grafana.GetAllDatasources(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	if cmd.verbose {
		fmt.Printf("Found %d datasources.\n", len(allDatasources))
	}
	for _, ds := range allDatasources {
		select {
		case <-cancel:
			exitBySignal()
		default:
			if datasources != nil {
				if _, ok := datasources[ds.Name]; !ok {
					continue
				}
			}
			if rawDs, err = json.Marshal(ds); err != nil {
				fmt.Fprintf(os.Stderr, "datasource marshal error %s\n", err)
				continue
			}
			var fname = fmt.Sprintf("%s.ds.%d.json", slug.Make(ds.Name), ds.OrgID)
			if err = ioutil.WriteFile(fname, rawDs, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, ds.Name))
				continue
			}
			if cmd.verbose {
				fmt.Printf("%s written into %s", ds.Name, fname)
			}
		}
	}
}

func extractDatasources(datasources map[string]bool, board sdk.Board) {
	for _, row := range board.Rows {
		for _, panel := range row.Panels {
			if panel.Datasource != nil {
				datasources[*panel.Datasource] = true
				fmt.Println(slug.Make(*panel.Datasource))
			}
		}
	}
}
