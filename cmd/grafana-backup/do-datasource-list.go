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

	"github.com/grafana-tools/sdk"
)

func doDatasourceList(opts ...option) {
	var (
		cmd         = initCommand(opts...)
		datasources []sdk.Datasource
		err         error
	)
	if datasources, err = cmd.grafana.GetAllDatasources(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	for _, ds := range datasources {
		select {
		case <-cancel:
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
		default:
			fmt.Printf("<%d> \"%s\" (%s) %s\n", ds.ID, ds.Name, ds.Type, ds.URL)
		}
		if cmd.verbose {
			fmt.Printf("Found %d datasources.\n", len(datasources))
		}
	}
Exit:
	fmt.Println()
}
