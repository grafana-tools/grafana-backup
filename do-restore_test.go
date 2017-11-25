package main


import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

type testresults struct {
	serverInstanceCalled *bool
	applyForCalled       *bool
	matchFilenameCalled  *bool
}

func newTestresults() testresults {

	value := retFalse()
	result := testresults{serverInstanceCalled: value, applyForCalled: value, matchFilenameCalled: value}
	result.matchFilenameCalled = value

	return result

	//var results testresults
	//results.serverInstanceCalled = false
	//results.applyForCalled       = &false
	//results.matchFilenameCalled  = &false
	//return results
}

func retFalse() *bool {
	fals := false
	return &fals
}

func retTrue() *bool {
	tru := true
	return &tru
}

func TestRestoreDashboards(t *testing.T) {
	t.Log("TestRestoreDashboards not yet implemented!")
}


//TODO: Create multiple tests which test things like sending multiple files
func TestRestoreDatasources(t *testing.T) {

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


	argCommand = "restore"
	argPath = "testdata/prometheus-test.ds.1.json"

	// For developing tests. Both of these cause this test to fail.
	//argPath = "testdata/*.1.json"
	//argPath = "testdata/promartheus-test.ds.1.json"

	// Some variables to track the results of the test

	// Check the accept header.
	acceptCorrect    := false
	// Check for some expected text in the post body.
	bodyCorrect      := false
	// Track how many times the API was called.
	numRequests      := 0
	// Were any requests made to other URIs?
	wrongUriRequests := false

	// Set up httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//TODO: Break this up into multiple functions so that the NoResponder doesn't cause us to fail Accept Header, body, etc.
	// Create a responder which will respond with valid JSON and check what was posted to us for validity.
	httpmock.RegisterResponder("POST", *flagServerURL + "/api/datasources",
		func(req *http.Request) (*http.Response, error) {

			numRequests++

			if strings.Contains(req.Header.Get("Accept"), "application/json") {
				acceptCorrect = true
			}

			//TODO: Expand this to unmarshal the JSON and check specific fields for specific values.

			// Get a string out of the io.ReadCloser
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			postBody := buf.String() // Does a complete copy of the bytes in the buffer.

			if strings.Contains(postBody, "prometheus-test") {
				bodyCorrect = true
			}

			// Uncomment for troubleshooting.
			//fmt.Printf("Request headers: \n%v\n", req.Header)
			//
			//fmt.Printf("Request body: \n%s\n", postBody)

			return httpmock.NewStringResponse(409, `{ "message": "This response is from the mocking framework!" }`), nil

			//TODO: Figure out how to make sure that do-restore is throwing an error when we return anything other than "Datasource added"
		},
	)

	httpmock.RegisterNoResponder(
		func(req *http.Request) (*http.Response, error) {

			wrongUriRequests = true

			fmt.Printf("Unexpected Request: \n%v\n", req)

			//fmt.Printf("Request headers: \n%v\n", req.Header)
			//
			//fmt.Printf("Request body: \n%v\n", req.Body)

			return httpmock.NewStringResponse(500, `{ "message": "Unexpected request" }`), nil
		},
	)

	doRestore(serverInstance, applyFor, matchFilename)

	if acceptCorrect != true {
		t.Error("Accept header was invalid.")
		//t.Fail()
	}

	if bodyCorrect != true {
		t.Error("Expected text not found in the POST body.")
		//t.Fail()
	}

	if numRequests != 1 {
		t.Errorf("The /api/datasources URI was called an incorrect number of times. Actual requests %d", numRequests)
	}

	if wrongUriRequests != false {
		t.Error("Request made to an unexpected URI. See the log for details.")
	}
}

//TODO: Change t.Log to t.Error when ready to implement this.
func TestRestoreUsers(t *testing.T) {
	t.Log("Test Restore Users not yet implemented because restoring users is not yet implemented.")
}
