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
	"io/ioutil"
	"os"

	"github.com/grafana-tools/sdk"
)

func doBackup(opts ...option) {
	var (
		cmd = initCommand(opts...)
	)
	if cmd.applyForHierarchy {
		//		backupDashboardHierarchy()
		backupDashboards(cmd)
		return
	}
	if cmd.applyForBoards {
		backupDashboards(cmd)
	}
	if cmd.applyForDs {
		//		backupDatasources()
	}
	if cmd.applyForUsers {
		//		backupUsers()
	}

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
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
		default:
			if rawBoard, meta, err = cmd.grafana.GetRawDashboard(link.URI); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, link.URI))
				continue
			}
			if err = ioutil.WriteFile(fmt.Sprintf("%s.json", meta.Slug), rawBoard, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, meta.Slug))
				continue
			}
			if cmd.verbose {
				fmt.Printf("%s.json backuped ok.\n", meta.Slug)
			}
		}
	}
Exit:
	fmt.Println()

}
