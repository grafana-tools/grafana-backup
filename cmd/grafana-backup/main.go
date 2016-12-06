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
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grafov/autograf/client"
)

type command struct {
	grafana   *client.Instance
	boardName string
	tags      []string
	starred   bool
	filenames string
}

type option func(*command) error

var (
	flagServerURL, flagServerKey          string
	flagTags, flagBoardName, flagFileName string
	flagTimeout                           time.Duration
	flagStarred, verbose                  bool
)

// TODO use first $XDG_CONFIG_HOME then try $XDG_CONFIG_DIRS
var tryConfigDirs = []string{"~/.config/grafana+", ".grafana+"}

func main() {
	// TODO parse config here

	flag.BoolVar(&verbose, "v", false, "verbose output")
	// Connection flags for single or two Grafana instances:
	flag.StringVar(&flagServerURL, "url", "", "URL of Grafana server")
	flag.StringVar(&flagServerKey, "key", "", "API key of Grafana server")
	flag.DurationVar(&flagTimeout, "timeout", 6*time.Minute, "read flagTimeout for interacting with Grafana (seconds)")
	// Dashboard matching flags:
	flag.StringVar(&flagTags, "tag", "", "dashboard should match all these tags")
	flag.BoolVar(&flagStarred, "starred", false, "only match starred dashboards")
	flag.StringVar(&flagBoardName, "name", "", "dashboard should match name")
	flag.StringVar(&flagFileName, "file", "", "use only listed files (file masks allowed)")
	flag.Parse()
	var args = flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	switch args[0] {
	case "backup":
		doBackup(serverInstance(), matchDashboard())
	case "restore":
		doRestore(serverInstance(), matchDashboard())
		// TBD
	case "ls", "list":
		// TBD
		// doList(matchDashboard())
	case "info":
		// TBD
		// doInfo(matchDashboard())
	case "config":
		// TBD
		// doConfig()
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Unknown command: %s\n\n", args[0]))
		printUsage()
		os.Exit(1)
	}
}

func serverInstance() option {
	return func(c *command) error {
		if flagServerURL != "" {
			return errors.New("you should provide the server URL")
		}
		if flagServerKey != "" {
			return errors.New("you should provide the server API key")
		}
		c.grafana = client.New(flagServerURL, flagServerKey, &http.Client{Timeout: flagTimeout})
		return nil
	}
}

func matchDashboard() option {
	return func(c *command) error {
		c.boardName = flagBoardName
		c.starred = flagStarred
		if flagTags != "" {
			for _, tag := range strings.Split(flagTags, ",") {
				c.tags = append(c.tags, strings.TrimSpace(tag))
			}
		}
		return nil
	}
}

func initCommand(opts ...option) *command {
	var (
		cmd = &command{}
		err error
	)
	for _, opt := range opts {
		if err = opt(cmd); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n\n", err))
			printUsage()
			os.Exit(1)
		}
	}
	return cmd
}

func printUsage() {
	fmt.Println(`Backup tool for Grafana.
Copyright (C) 2016  Alexander I.Grafov <siberian@laika.name>

This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under conditions of GNU GPL license v3.

Usage: $ grafana-backup [flags] <command>
`)
	flag.PrintDefaults()

}
