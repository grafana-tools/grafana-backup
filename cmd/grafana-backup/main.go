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
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grafov/autograf/client"
)

func main() {
	var (
		serverURL, serverKey, tagline, boardName, fileMask string
		starred, verbose                                   bool
	)
	flag.StringVar(&serverURL, "url", "http://localhost:3000", "URL of Grafana server")
	flag.StringVar(&serverKey, "key", "", "API key of Grafana server")
	flag.StringVar(&tagline, "tag", "", "dashboard should match all these tags")
	flag.BoolVar(&starred, "starred", false, "only match starred dashboards")
	flag.StringVar(&boardName, "name", "", "dashboard should match name")
	flag.StringVar(&fileMask, "file", "", "use only listed files (file masks allowed)")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()
	var args = flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	var tags []string
	if tagline != "" {
		for _, tag := range strings.Split(tagline, ",") {
			tags = append(tags, tag)
		}
	}

	// TODO parse config here

	switch args[0] {
	case "backup":
		doBackup(serverURL, serverKey, starred, boardName, tags, verbose)
	case "restore":
		// TBD
	case "ls", "list":
		// TBD
	case "info":
		// TBD
	case "config":
		// TBD
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Unknown command: %s\n\n", args[0]))
		printUsage()
		os.Exit(1)
	}
}

func doBackup(serverURL, serverKey string, starred bool, boardName string, tags []string, verbose bool) {
	var (
		boardLinks []client.FoundBoard
		rawBoard   []byte
		meta       client.BoardProperties
		err        error
	)
	c := client.New(serverURL, serverKey, &http.Client{Timeout: 6 * time.Minute})
	if boardLinks, err = c.SearchDashboards(boardName, starred, tags...); err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	if verbose {
		fmt.Printf("Found %d dashboards that matched the conditions.\n", len(boardLinks))
	}
	var cancel = make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
	for _, link := range boardLinks {
		select {
		case <-cancel:
			fmt.Fprintf(os.Stderr, "Execution was cancelled.")
			goto Exit
		default:
			if rawBoard, meta, err = c.GetRawDashboard(link.URI); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, link.URI))
				continue
			}
			if err = ioutil.WriteFile(fmt.Sprintf("%s.json", meta.Slug), rawBoard, os.FileMode(int(0666))); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%s for %s\n", err, meta.Slug))
				continue
			}
			if verbose {
				fmt.Printf("%s.json backuped ok.\n", meta.Slug)
			}
		}
	}
Exit:
	fmt.Println()
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
