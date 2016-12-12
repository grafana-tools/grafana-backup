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
	"os"

	"github.com/grafov/autograf/client"
)

func doDashboardList(opts ...option) {
	var (
		cmd         = initCommand(opts...)
		foundBoards []client.FoundBoard
		err         error
	)
	if foundBoards, err = cmd.grafana.SearchDashboards(cmd.boardTitle, cmd.starred, cmd.tags...); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	for _, meta := range foundBoards {
		select {
		case <-cancel:
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
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
Exit:
	fmt.Println()
}
