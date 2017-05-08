package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type header struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type testInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Version     string `json:"version"`
	DateUpdated string `json:"date_uploaded"`
	Author      string `json:"author"`
}

type payload struct {
	Headers []header `json:"headers"`
	Body    string   `json:"body"`
}

type request struct {
	Verb    string  `json:"verb"`
	URL     string  `json:"url"`
	Payload payload `json:"payload"`
}

type expect struct {
	ParseAs      string      `json:"parse_as"`
	HTTPCode     int64       `json:"http_code"`
	MaxLatencyMS int64       `json:"max_latency_ms"`
	Headers      []header    `json:"headers"`
	Body         interface{} `json:"body"`
}

type actual struct {
	HTTPCode  int             `json:"http_code"`
	LatencyMS int64           `json:"latency_ms"`
	Headers   json.RawMessage `json:"headers"`
	Body      json.RawMessage `json:"body"`
}

type testResult struct {
	PassFail  string   `json:"pass_fail"`
	Timestamp string   `json:"timestamp"`
	TestInfo  testInfo `json:"test_info"`
	Request   request  `json:"request"`
	Expect    expect   `json:"expect"`
	Actual    actual   `json:"actual"`
}

type testRequest struct {
	TestInfo testInfo `json:"test_info"`
	Request  request  `json:"request"`
	Expect   expect   `json:"expect"`
}

// readTestJSON reads JSON input from stdin and returns it as a formatted Go
// struct
func readTestJSON() testRequest {
	j, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("Error reading content from stdin")
		panic(err)
	}
	var tr testRequest
	err = json.Unmarshal(j, &tr)
	if err != nil {
		log.Println("Error parsing content read from stdin as JSON")
		log.Printf("%v\n", string(j))
		panic(err)
	}

	return tr
}

func populateRequest(testCaseRequest testRequest) (testInfo, request, expect) {
	testinfo := &testInfo{
		ID:          testCaseRequest.TestInfo.ID,
		Description: testCaseRequest.TestInfo.Description,
		Version:     testCaseRequest.TestInfo.Version,
		DateUpdated: testCaseRequest.TestInfo.DateUpdated,
		Author:      testCaseRequest.TestInfo.Author,
	}

	request := &request{
		Verb: testCaseRequest.Request.Verb,
		URL:  testCaseRequest.Request.URL,
	}

	expect := &expect{
		ParseAs:      testCaseRequest.Expect.ParseAs,
		HTTPCode:     testCaseRequest.Expect.HTTPCode,
		MaxLatencyMS: testCaseRequest.Expect.MaxLatencyMS,
		Headers:      testCaseRequest.Expect.Headers,
		Body:         testCaseRequest.Expect.Body,
	}
	return *testinfo, *request, *expect
}

// executeRequest executes the JSON request defined in the test case, and captures & returns
// the response body, response headers, HTTP status and latency
func executeRequest(request request) (interface{}, interface{}, int, time.Duration) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(request.Verb, request.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	startTime := time.Now()
	resp, err := httpClient.Do(req)
	endTime := time.Since(startTime)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	responseDecoder := json.NewDecoder(resp.Body)
	var v interface{} // Not sure what the response will look like, so just implement an interface
	err = responseDecoder.Decode(&v)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Printf("v\n%v\n", v)
	headers := resp.Header
	httpCode := resp.StatusCode
	latency := endTime
	return v, headers, httpCode, latency
}

func populateResponse(body interface{}, headers interface{}, httpCode int, latency time.Duration) actual {
	var actual actual
	actual.HTTPCode = httpCode
	actual.LatencyMS = int64(latency / time.Millisecond)

	var bodyStr json.RawMessage
	bodyStr, err := json.Marshal(body)
	if err != nil {
		panic("Unable to parse response body as JSON")
	}
	//log.Printf("bodyStr\n%v\n", string(bodyStr))
	actual.Body = bodyStr

	var headerStr json.RawMessage
	headerStr, err = json.Marshal(headers)
	if err != nil {
		panic("Unable to parse response headers as JSON")
	}
	actual.Headers = headerStr

	return actual
}

// compareActualVersusExpected compares the actual response against the
// expected response, and returns a boolean indicating whether the match was
// good or bad
func compareActualVersusExpected(actual actual, expect expect) bool {

	switch expect.ParseAs {
	case "":
		// if there's no parser defined, return false; this is a viable approach
		// for simply collecting info about an API response but not a test case
		return false
	case "regex":
		// we want to parse the actual content against regex patterns that are defined
		// in the "expected" part of the test case
		return true
	case "exact_match":
		// we want the actual response to be an exact match to the "expected" part of
		// the test case. If there are extra fields in the actual response to what's
		// defined in "expected", then the test should fail
		return false
	case "partial_match":
		// we want the actual response fields to be an exact match to the "expected" fields
		// defined in the test case, but the "expected" fields may not contain all the fields in
		// the actual response
		return false
	default:
		panic("expect.parse_as should be one of 'regex', 'exact_match', 'partial_match'")
	}
}

func main() {

	// first read the JSON test request from stdin
	testCaseRequest := readTestJSON()

	// then populate the the "request" content that will eventually be sent to stdout
	testinfo, request, expect := populateRequest(testCaseRequest)

	// then execute the request, and capture the response body, headers, http status and latency
	body, headers, httpCode, latency := executeRequest(request)

	// then populate the "response" content that will eventually be sent to stdout
	actual := populateResponse(body, headers, httpCode, latency)

	// then check whether the "response" matches what was expected
	var passFail = "fail"
	if compareActualVersusExpected(actual, expect) {
		passFail = "pass"
	}

	testresult := &testResult{
		PassFail:  passFail,
		Timestamp: time.Now().Local().Format(time.RFC3339),
		Request:   request,
		TestInfo:  testinfo,
		Expect:    expect,
		Actual:    actual,
	}

	testresultJSON, err := json.MarshalIndent(testresult, "", "  ")
	if err != nil {
		panic("Unable to display output as JSON")
	}
	fmt.Printf("%+v\n", string(testresultJSON))
}
