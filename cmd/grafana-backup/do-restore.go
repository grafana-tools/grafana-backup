package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func doRestore(opts ...option) {
	var (
		cmd        *command
		filesInDir []os.FileInfo
		err        error
	)
	cmd = initCommand(opts...)
	filesInDir, err = ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	var rawBoard []byte
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawBoard, err = ioutil.ReadFile(file.Name()); err != nil {
				fmt.Fprintf(os.Stderr, "error on read %s", file.Name())
				continue
			}
			if err = cmd.grafana.SetRawDashboard(rawBoard); err != nil {
				fmt.Fprintf(os.Stderr, "error on importing dashboard from %s", file.Name())
				continue
			}
		}
	}

}
