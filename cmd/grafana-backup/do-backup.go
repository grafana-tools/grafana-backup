package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/grafov/autograf/client"
)

func doBackup(opts ...option) {
	cmd := initCommand(opts...)
	var (
		boardLinks []client.FoundBoard
		rawBoard   []byte
		meta       client.BoardProperties
		err        error
	)
	if boardLinks, err = cmd.grafana.SearchDashboards(cmd.boardName, cmd.starred, cmd.tags...); err != nil {
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
			if rawBoard, meta, err = cmd.grafana.GetRawDashboard(link.URI); err != nil {
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
