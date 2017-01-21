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
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/grafana-tools/sdk"
)

var (
	// Connection flags.
	flagServerURL = flag.String("url", "", "URL of Grafana server")
	flagServerKey = flag.String("key", "", "API key of Grafana server")
	flagTimeout   = flag.Duration("timeout", 6*time.Minute, "read flagTimeout for interacting with Grafana in seconds")

	// Dashboard matching flags.
	flagTags       = flag.String("tag", "", "dashboard should match all these tags")
	flagBoardTitle = flag.String("title", "", "dashboard title should match name")
	flagStarred    = flag.Bool("starred", false, "only match starred dashboards")
	// Common flags.
	flagApplyFor = flag.String("apply-for", "auto", `apply operation only for some kind of objects, available values are "auto", "all", "dashboards", "datasources", "users"`)
	flagForce    = flag.Bool("force", false, "force overwrite of existing objects")
	flagVerbose  = flag.Bool("verbose", false, "verbose output")

	// The args after flags.
	argCommand string
	argPath    = "*"
)

var cancel = make(chan os.Signal, 1)

// TODO use first $XDG_CONFIG_HOME then try $XDG_CONFIG_DIRS
var tryConfigDirs = []string{"~/.config/grafana+", ".grafana+"}

func main() {
	// TODO parse config here

	flag.Parse()
	if flag.NArg() == 0 {
		printUsage()
		os.Exit(2)
	}
	var args = flag.Args()
	// First mandatory argument is command.
	argCommand = args[0]
	// Second optional argument is file path.
	if flag.NArg() > 1 {
		argPath = args[1]
	}
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
	switch argCommand {
	case "backup":
		// TODO fix logic accordingly with apply-for
		doBackup(serverInstance, applyFor, matchDashboard)
	case "restore":
		// TODO fix logic accordingly with apply-for
		doRestore(serverInstance, applyFor, matchFilename)
	case "ls":
		// TODO fix logic accordingly with apply-for
		doObjectList(serverInstance, applyFor, matchDashboard)
	case "ls-files":
		// TODO merge this command with ls
		doFileList(matchFilename, applyFor, matchDashboard)
	case "config-set":
		// TBD
		// doConfigSet()
	case "config-get":
		// TBD
		// doConfigGet()
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("unknown command: %s\n\n", args[0]))
		printUsage()
		os.Exit(1)
	}
}

type command struct {
	grafana             *sdk.Client
	applyHierarchically bool
	applyForBoards      bool
	applyForDs          bool
	applyForUsers       bool
	boardTitle          string
	tags                []string
	starred             bool
	filenames           []string
	force               bool
	verbose             bool
}

type option func(*command) error

func serverInstance(c *command) error {
	if *flagServerURL == "" {
		return errors.New("you should provide the server URL")
	}
	if *flagServerKey == "" {
		return errors.New("you should provide the server API key")
	}
	c.grafana = sdk.NewClient(*flagServerURL, *flagServerKey, &http.Client{Timeout: *flagTimeout})
	return nil
}

func applyFor(c *command) error {
	if *flagApplyFor == "" {
		return fmt.Errorf("flag '-apply-for' provided with empty argument")
	}
	for _, objectKind := range strings.Split(strings.ToLower(*flagApplyFor), ",") {
		switch objectKind {
		case "auto":
			c.applyForHierarchy = true
			c.applyForBoards = true
			c.applyForDs = true
		case "all":
			c.applyForBoards = true
			c.applyForDs = true
			c.applyForUsers = true
		case "dashboards":
			c.applyForBoards = true
		case "datasources":
			c.applyForDs = true
		case "users":
			c.applyForUsers = true
		default:
			return fmt.Errorf("unknown argument %s", objectKind)
		}
	}
	return nil
}

func matchDashboard(c *command) error {
	c.boardTitle = *flagBoardTitle
	c.starred = *flagStarred
	if *flagTags != "" {
		for _, tag := range strings.Split(*flagTags, ",") {
			c.tags = append(c.tags, strings.TrimSpace(tag))
		}
	}
	return nil
}

func matchFilename(c *command) error {
	var (
		files []string
		err   error
	)
	if files, err = filepath.Glob(argPath); err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("there are no files matching selected pattern found")
	}
	c.filenames = files
	return nil
}

func initCommand(opts ...option) *command {
	var (
		cmd = &command{force: *flagForce, verbose: *flagVerbose}
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
Copyright (C) 2016-2017  Alexander I.Grafov <siberian@laika.name>

This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under conditions of GNU GPL license v3.

Usage: $ grafana-backup [flags] <command>

Available commands are: backup, restore, list, info, config, help.
Call 'grafana-backup help <command>' for details about the command.
`)
	flag.PrintDefaults()

}

func exitBySignal() {
	fmt.Fprintf(os.Stderr, "Execution was cancelled.\n")
	os.Exit(1)
}
