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
	"fmt"
	"os"

	"github.com/grafana-tools/sdk"
)

func doObjectList(opts ...option) {
	var (
		cmd = initCommand(opts...)
		err error
	)
	if cmd.applyForBoards {
		listDashboards(cmd)
	}
	if cmd.applyForDs {
		listDatasources(cmd)
	}
	if cmd.applyForUsers {
		listUsers(cmd)
	}
}

func listDashboards(cmd *command) {
	var (
		foundBoards []sdk.FoundBoard
		err         error
	)
	if foundBoards, err = cmd.grafana.SearchDashboards(cmd.boardTitle, cmd.starred, cmd.tags...); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		return
	}
	for _, meta := range foundBoards {
		select {
		case <-cancel:
			exit()
		default:
			fmt.Printf("<%d> \"%s\" %v ", meta.ID, meta.Title, meta.Tags)
			if meta.IsStarred {
				fmt.Print("*")
			}
			fmt.Println()
		}
		if cmd.verbose {
			fmt.Printf("Found %d dashboards.\n", len(foundBoards))
		}
	}
}

func listDatasources(cmd *command) {
	var (
		datasources []sdk.Datasource
		err         error
	)
	if datasources, err = cmd.grafana.GetAllDatasources(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		return
	}
	for _, ds := range datasources {
		select {
		case <-cancel:
			exit()
		default:
			fmt.Printf("<%d> \"%s\" (%s) %s\n", ds.ID, ds.Name, ds.Type, ds.URL)
		}
		if cmd.verbose {
			fmt.Printf("Found %d datasources.\n", len(datasources))
		}
	}
}

func listUsers(cmd *command) {
	var (
		allUsers []sdk.User
		err      error
	)
	if allUsers, err = cmd.grafana.GetAllUsers(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		return
	}
	for _, user := range allUsers {
		select {
		case <-cancel:
			exit()
		default:
			fmt.Printf("%s \"%s\" <%s>", user.Login, user.Name, user.Email)
			if user.IsGrafanaAdmin {
				fmt.Print(" admin")
			}
			fmt.Println()
		}
		if cmd.verbose {
			fmt.Printf("Found %d users.\n", len(users))
		}
	}
}
