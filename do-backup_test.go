package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"os"
)

func TestBackupDashboards(t *testing.T) {
	*flagServerURL = "http://noserver.nodomain.com:3123"
	*flagServerKey = "thisisnotreallyanapikey"
	*flagVerbose  = true
	*flagApplyFor = "dashboards"
	*flagDir      = "/var/tmp/testDashboardsOutDir"
	argCommand = "backup"

	// Some variables to track the results of the test

	// Check the accept header.
	acceptCorrect    := false
	// Check the content type.
	cTypeCorrect    := false
	// Track how many times the API was called.
	numRequests      := 0
	// Were any requests made to other URIs?
	wrongUriRequests := false

	// Set up httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//TODO: Break this up into multiple functions so that the NoResponder doesn't cause us to fail Accept Header, body, etc.
	// Create a responder for /api/search which resturns the list of dashboards (and a tiny bit more.)
	httpmock.RegisterResponder("GET", *flagServerURL + "/api/search",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			// This is the kind of stuff we see in a request. I don't think it's neccessary to check things like the host and URL but maybe.
			//Request headers:
			//map[Content-Type:[application/json] User-Agent:[autograf] Authorization:[Bearer thisisnotreallyanapikey] Accept:[application/json]]
			//Request Host: noserver.nodomain.com:3123
			//Request Method: GET
			//Request Protocol: HTTP/1.1
			//Request RemoteAddr:
			//Request URI:
			//Request URL: http://noserver.nodomain.com:3123/api/datasources

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			if strings.Contains(req.Header.Get("Content-type"), "application/json") {
				cTypeCorrect = true
			}

			return httpmock.NewStringResponse(200,
				`[
				  {
					"id": 31,
					"title": "Test Dashboard 1",
					"uri": "db/test-dashboard-1",
					"type": "dash-db",
					"tags": [
					  "tagone",
					  "tagtwo"
					],
					"isStarred": false
				  },
				  {
					"id": 32,
					"title": "Test Dashboard 2",
					"uri": "db/test-dashboard-2",
					"type": "dash-db",
					"tags": [],
					"isStarred": false
				  },
				  {
					"id": 33,
					"title": "Test Dashboard 3",
					"uri": "db/test-dashboard-3",
					"type": "dash-db",
					"tags": [],
					"isStarred": false
				  }
				]
`), nil

		},
	)

	// TODO: The requests for these 3 dashboards are being made but the files are not being created. There must be something wrong with the syntax.
	// TODO: Compact these up a bit so that they aren't taking up so much screen real-estate
	// Create a responder for Test Dashboard 1
	httpmock.RegisterResponder("GET", *flagServerURL + "/api/dashboards/db/test-dashboard-1",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			if strings.Contains(req.Header.Get("Content-type"), "application/json") {
				cTypeCorrect = true
			}

			return httpmock.NewStringResponse(200,
				`{"meta":{"type":"db","canSave":true,"canEdit":true,"canStar":true,"slug":"test-dashboard-1","expires":"0001-01-01T00:00:00Z","created":"2017-11-24T09:05:33-08:00","updated":"2017-11-24T09:35:26-08:00","updatedBy":"admin","createdBy":"admin","version":3},"dashboard":{"annotations":{"list":[]},"editable":true,"gnetId":null,"graphTooltip":0,"hideControls":false,"id":31,"links":[],"rows":[{"collapse":false,"height":250,"panels":[{"aliasColors":{},"bars":false,"dashLength":10,"dashes":false,"datasource":"prometheus-test","fill":1,"id":2,"legend":{"avg":false,"current":false,"max":false,"min":false,"show":true,"total":false,"values":false},"lines":true,"linewidth":1,"links":[],"nullPointMode":"null","percentage":false,"pointradius":5,"points":false,"renderer":"flot","seriesOverrides":[],"spaceLength":10,"span":12,"stack":false,"steppedLine":false,"targets":[{"expr":"( \t100  \t*  \t( \t\t1  \t\t-  \t\tavg( \t\t\tirate( \t\t\t\tnode_cpu{job=\"node\",mode=\"idle\",Team=~\"$Team\"}[5m] \t\t\t) \t\t) BY (instance) \t) )","format":"time_series","intervalFactor":2,"legendFormat":"{{ instance }}","refId":"A","step":20}],"thresholds":[],"timeFrom":null,"timeShift":null,"title":"CPU by Team","tooltip":{"shared":true,"sort":0,"value_type":"individual"},"type":"graph","xaxis":{"buckets":null,"mode":"time","name":null,"show":true,"values":[]},"yaxes":[{"format":"short","label":null,"logBase":1,"max":null,"min":null,"show":true},{"format":"short","label":null,"logBase":1,"max":null,"min":null,"show":true}]}],"repeat":null,"repeatIteration":null,"repeatRowId":null,"showTitle":false,"title":"Dashboard Row","titleSize":"h6"}],"schemaVersion":14,"style":"dark","tags":["tagone","tagtwo"],"templating":{"list":[{"allValue":null,"current":{"text":"ISLab","value":"ISLab"},"datasource":"prometheus-test","hide":0,"includeAll":false,"label":null,"multi":false,"name":"Team","options":[],"query":"label_values(node_boot_time, Team)","refresh":2,"regex":"","sort":0,"tagValuesQuery":"","tags":[],"tagsQuery":"","type":"query","useTags":false}]},"time":{"from":"now-6h","to":"now"},"timepicker":{"refresh_intervals":["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],"time_options":["5m","15m","1h","6h","12h","24h","2d","7d","30d"]},"timezone":"","title":"Test Dashboard 1","version":3}}`), nil

		},
	)

	// Create a responder for Test Dashboard 2
	httpmock.RegisterResponder("GET", *flagServerURL + "/api/dashboards/db/test-dashboard-2",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			if strings.Contains(req.Header.Get("Content-type"), "application/json") {
				cTypeCorrect = true
			}

			return httpmock.NewStringResponse(200,
				`{"meta":{"type":"db","canSave":true,"canEdit":true,"canStar":true,"slug":"test-dashboard-2","expires":"0001-01-01T00:00:00Z","created":"2017-11-24T09:36:03-08:00","updated":"2017-11-24T09:36:03-08:00","updatedBy":"admin","createdBy":"admin","version":1},"dashboard":{"annotations":{"list":[]},"editable":true,"gnetId":null,"graphTooltip":0,"hideControls":false,"id":32,"links":[],"rows":[{"collapse":false,"height":"250px","panels":[{"content":"# This is a Title","id":1,"links":[],"mode":"markdown","span":12,"title":"A Panel Title","type":"text"}],"repeat":null,"repeatIteration":null,"repeatRowId":null,"showTitle":false,"title":"Dashboard Row","titleSize":"h6"}],"schemaVersion":14,"style":"dark","tags":[],"templating":{"list":[]},"time":{"from":"now-6h","to":"now"},"timepicker":{"refresh_intervals":["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],"time_options":["5m","15m","1h","6h","12h","24h","2d","7d","30d"]},"timezone":"","title":"Test Dashboard 2","version":1}}`), nil

		},
	)

	// Create a responder for Test Dashboard 3
	httpmock.RegisterResponder("GET", *flagServerURL + "/api/dashboards/db/test-dashboard-3",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			if strings.Contains(req.Header.Get("Content-type"), "application/json") {
				cTypeCorrect = true
			}

			return httpmock.NewStringResponse(200,
				`{"meta":{"type":"db","canSave":true,"canEdit":true,"canStar":true,"slug":"test-dashboard-3","expires":"0001-01-01T00:00:00Z","created":"2017-11-24T09:36:44-08:00","updated":"2017-11-24T09:38:25-08:00","updatedBy":"admin","createdBy":"admin","version":2},"dashboard":{"annotations":{"list":[]},"editable":true,"gnetId":null,"graphTooltip":0,"hideControls":false,"id":33,"links":[],"rows":[{"collapse":false,"height":"250px","panels":[{"content":"# This is also a title","id":1,"links":[],"mode":"markdown","span":12,"title":"Still the Panel Title","type":"text"}],"repeat":null,"repeatIteration":null,"repeatRowId":null,"showTitle":false,"title":"Dashboard Row","titleSize":"h6"},{"collapse":false,"height":250,"panels":[{"cacheTimeout":null,"colorBackground":false,"colorValue":false,"colors":["rgba(245, 54, 54, 0.9)","rgba(237, 129, 40, 0.89)","rgba(50, 172, 45, 0.97)"],"datasource":"Promt2Local","format":"none","gauge":{"maxValue":100,"minValue":0,"show":false,"thresholdLabels":false,"thresholdMarkers":true},"id":2,"interval":null,"links":[],"mappingType":1,"mappingTypes":[{"name":"value to text","value":1},{"name":"range to text","value":2}],"maxDataPoints":100,"nullPointMode":"connected","nullText":null,"postfix":"","postfixFontSize":"50%","prefix":"","prefixFontSize":"50%","rangeMaps":[{"from":"null","text":"N/A","to":"null"}],"span":12,"sparkline":{"fillColor":"rgba(31, 118, 189, 0.18)","full":false,"lineColor":"rgb(31, 120, 193)","show":false},"tableColumn":"","targets":[{"expr":"( \t100  \t*  \t( \t\t1  \t\t-  \t\tavg( \t\t\tirate( \t\t\t\tnode_cpu{job=\"node\",mode=\"idle\",Team=~\"$Team\"}[5m] \t\t\t) \t\t) BY (instance) \t) )","format":"time_series","intervalFactor":2,"legendFormat":"{{ instance }}","refId":"A","step":600}],"thresholds":"","title":"Singlestat Panel Title","type":"singlestat","valueFontSize":"80%","valueMaps":[{"op":"=","text":"N/A","value":"null"}],"valueName":"avg"}],"repeat":null,"repeatIteration":null,"repeatRowId":null,"showTitle":false,"title":"Dashboard Row","titleSize":"h6"}],"schemaVersion":14,"style":"dark","tags":[],"templating":{"list":[]},"time":{"from":"now-6h","to":"now"},"timepicker":{"refresh_intervals":["5s","10s","30s","1m","5m","15m","30m","1h","2h","1d"],"time_options":["5m","15m","1h","6h","12h","24h","2d","7d","30d"]},"timezone":"","title":"Test Dashboard 3","version":2}}`), nil

		},
	)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {

			wrongUriRequests = true

			// Uncomment for troubleshooting.
			fmt.Printf("Request headers: \n%v\n", req.Header)
			fmt.Printf("Request Host: %s\n", req.Host)
			fmt.Printf("Request Method: %s\n", req.Method)
			fmt.Printf("Request Protocol: %s\n", req.Proto)
			fmt.Printf("Request RemoteAddr: %s\n", req.RemoteAddr)
			fmt.Printf("Request URI: %s\n", req.RequestURI)
			fmt.Printf("Request URL: %s\n", req.URL)


			// Uncomment for moar troubleshooting
			//// Get a string out of the io.ReadCloser
			//buf := new(bytes.Buffer)
			//buf.ReadFrom(req.Body)
			//reqBody := buf.String() // Does a complete copy of the bytes in the buffer.
			//
			//fmt.Printf("Request body: \n%s\n", reqBody)


			// Uncomment for the mostest troubleshooting
			fmt.Printf("Unexpected Request: \n%s\n", req)

			return httpmock.NewStringResponse(500, `{ "message": "Unexpected request" }`), nil
		},
	)

	// Trigger the actual backup
	doBackup(serverInstance, applyFor, matchDashboard)


	if acceptCorrect != true {
		t.Error("Accept header was invalid.")
		//t.Fail()
	}

	// There should be 4 requests. One for search and one for each of the 3 dashboards.
	if numRequests != 4 {
		t.Errorf("The /api/search URI was called an incorrect number of times. Actual requests %d", numRequests)
	}

	if wrongUriRequests != false {
		t.Error("Request made to an unexpected URI. See the log for details.")
	}

	// Check the output files and confirm that they are correct.
	//TODO: Confirm that they are correct instead of just confirming that they exist.

	//name := "FileOrDir"
	fi, err := os.Stat("/var/tmp/testDashboardsOutDir")
	if err != nil {
		fmt.Println(err)
		return
	}

	if !fi.Mode().IsDir() {
		t.Error("Output directory does not exist or is not a directory.")
	}

	if _, err := os.Stat("/var/tmp/testDashboardsOutDir/test-dashboard-1.db.json"); os.IsNotExist(err) {
		t.Error("test-dashboard-1 output file does not exist.")
	}

	if _, err := os.Stat("/var/tmp/testDashboardsOutDir/test-dashboard-2.db.json"); os.IsNotExist(err) {
		t.Error("test-dashboard-2 output file does not exist.")
	}

	if _, err := os.Stat("/var/tmp/testDashboardsOutDir/test-dashboard-3.db.json"); os.IsNotExist(err) {
		t.Error("test-dashboard-3 output file does not exist.")
	}

	// Cleanup the output files. Clean up each file and directory explicitly so that we don't accidentally rm -rf something important.

	//err = os.Remove("/var/tmp/testOutputDir/promt2local.ds.1.json")
	//err = os.Remove("/var/tmp/testOutputDir/prometheus-test.ds.1.json")
	//err = os.Remove("/var/tmp/testOutputDir")
	//
	//if err != nil {
	//	t.Errorf("Unable to clean up output file %s", err)
	//}


}


//TODO: Create multiple tests which test things like sending multiple files
func TestBackupDatasources(t *testing.T) {

	//flagServerURL = flag.String("url", "", "URL of Grafana server")
	*flagServerURL = "http://noserver.nodomain.com:3123"
	//flagServerKey = flag.String("key", "", "API key of Grafana server")
	*flagServerKey = "thisisnotreallyanapikey"
	//flagTimeout   = flag.Duration("timeout", 6*time.Minute, "read flagTimeout for interacting with Grafana in seconds")

	//// Dashboard matching flags.
	//flagTags       = flag.String("tag", "", "dashboard should match all these tags")
	//flagBoardTitle = flag.String("title", "", "dashboard title should match name")
	//flagStarred    = flag.Bool("starred", false, "only match starred dashboards")

	//// Common flags.
	//flagApplyFor = flag.String("apply-for", "auto", `apply operation only for some kind of objects, available values are "auto", "all", "dashboards", "datasources", "users"`)
	*flagApplyFor = "datasources"
	//flagForce    = flag.Bool("force", false, "force overwrite of existing objects")
	//flagVerbose  = flag.Bool("verbose", false, "verbose output")
	//flagDir      = flag.String("dir", "backup", "A directory to write backup files to or read them from.")
	*flagDir      = "/var/tmp/testOutputDir"

	argCommand = "backup"
	// These were used in the test restore datasources
	//argPath = "testdata/prometheus-test.ds.1.json"
	//argPath = "testdata/*.1.json"
	//argPath = "testdata/promartheus-test.ds.1.json"

	// Some variables to track the results of the test

	// Check the accept header.
	acceptCorrect    := false
	// Check the content type.
	cTypeCorrect    := false
	// Track how many times the API was called.
	numRequests      := 0
	// Were any requests made to other URIs?
	wrongUriRequests := false

	// Set up httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//TODO: Break this up into multiple functions so that the NoResponder doesn't cause us to fail Accept Header, body, etc.
	// Create a responder which will respond with valid JSON and check what was posted to us for validity.
	httpmock.RegisterResponder("GET", *flagServerURL + "/api/datasources",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			// This is the kind of stuff we see in a request. I don't think it's neccessary to check things like the host and URL but maybe.
			//Request headers:
			//map[Content-Type:[application/json] User-Agent:[autograf] Authorization:[Bearer thisisnotreallyanapikey] Accept:[application/json]]
			//Request Host: noserver.nodomain.com:3123
			//Request Method: GET
			//Request Protocol: HTTP/1.1
			//Request RemoteAddr:
			//Request URI:
			//Request URL: http://noserver.nodomain.com:3123/api/datasources

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			if strings.Contains(req.Header.Get("Content-type"), "application/json") {
				cTypeCorrect = true
			}

			return httpmock.NewStringResponse(200,
			`[
			{
			"id": 10,
			"orgId": 1,
			"name": "prometheus-test",
			"type": "prometheus",
			"typeLogoUrl": "public/app/plugins/datasource/prometheus/img/prometheus_logo.svg",
			"access": "direct",
			"url": "http://prometheus-test.example.com:9090",
			"password": "",
			"user": "",
			"database": "",
			"basicAuth": false,
			"isDefault": false,
			"jsonData": {}
			},
			{
			"id": 8,
			"orgId": 1,
			"name": "Promt2Local",
			"type": "prometheus",
			"typeLogoUrl": "public/app/plugins/datasource/prometheus/img/prometheus_logo.svg",
			"access": "direct",
			"url": "http://prometheus-test.example.com:9090",
			"password": "",
			"user": "",
			"database": "",
			"basicAuth": false,
			"isDefault": false,
			"jsonData": {}
			}
		]`), nil

		},
	)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {

			wrongUriRequests = true

			// Uncomment for troubleshooting.
			fmt.Printf("Request headers: \n%v\n", req.Header)
			fmt.Printf("Request Host: %s\n", req.Host)
			fmt.Printf("Request Method: %s\n", req.Method)
			fmt.Printf("Request Protocol: %s\n", req.Proto)
			fmt.Printf("Request RemoteAddr: %s\n", req.RemoteAddr)
			fmt.Printf("Request URI: %s\n", req.RequestURI)
			fmt.Printf("Request URL: %s\n", req.URL)


			// Uncomment for moar troubleshooting
			//// Get a string out of the io.ReadCloser
			//buf := new(bytes.Buffer)
			//buf.ReadFrom(req.Body)
			//reqBody := buf.String() // Does a complete copy of the bytes in the buffer.
			//
			//fmt.Printf("Request body: \n%s\n", reqBody)


			// Uncomment for the mostest troubleshooting
			fmt.Printf("Unexpected Request: \n%s\n", req)

			return httpmock.NewStringResponse(500, `{ "message": "Unexpected request" }`), nil
		},
	)

	doBackup(serverInstance, applyFor, matchDashboard)

	if acceptCorrect != true {
		t.Error("Accept header was invalid.")
		//t.Fail()
	}

	if numRequests != 1 {
		t.Errorf("The /api/datasources URI was called an incorrect number of times. Actual requests %d", numRequests)
	}

	if wrongUriRequests != false {
		t.Error("Request made to an unexpected URI. See the log for details.")
	}

	// Check the output files and confirm that they are correct.
	//TODO: Confirm that they are correct instead of just confirming that they exist.

	//name := "FileOrDir"
	fi, err := os.Stat("/var/tmp/testOutputDir/")
	if err != nil {
		fmt.Println(err)
		return
	}

	if !fi.Mode().IsDir() {
		t.Error("Output directory does not exist or is not a directory.")
	}

	if _, err := os.Stat("/var/tmp/testOutputDir/prometheus-test.ds.1.json"); os.IsNotExist(err) {
		t.Error("prometheus-test output file does not exist.")
	}

	if _, err := os.Stat("/var/tmp/testOutputDir/promt2local.ds.1.json"); os.IsNotExist(err) {
		t.Error("prom2local file does not exist.")
	}

	// Cleanup the output files. Clean up each file and directory explicitly so that we don't accidentally rm -rf something important.

	err = os.Remove("/var/tmp/testOutputDir/promt2local.ds.1.json")
	err = os.Remove("/var/tmp/testOutputDir/prometheus-test.ds.1.json")
	err = os.Remove("/var/tmp/testOutputDir")

	if err != nil {
		t.Errorf("Unable to clean up output file %s", err)
	}


}

//TODO: Change t.Log to t.Error when ready to implement this.
func TestRestoreUsers(t *testing.T) {
	t.Log("Test Restore Users not yet implemented because restoring users is not yet implemented.")
}
