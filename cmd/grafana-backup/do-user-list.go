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

func doUserList(opts ...option) {
	var (
		cmd   = initCommand(opts...)
		users []sdk.User
		err   error
	)
	if users, err = cmd.grafana.GetAllUsers(); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	for _, user := range users {
		select {
		case <-cancel:
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
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
Exit:
	fmt.Println()
}
