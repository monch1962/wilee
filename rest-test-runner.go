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

/*type TestInput struct {
	testinfo TestInfo //`json:"test_info"`
}*/

func readTestJSON() testRequest {
	// Read JSON input from stdin and return as a formatted Go struct
	j, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("Error reading content from stdin")
		panic(err)
	}
	var tr testRequest
	err = json.Unmarshal(j, &tr)
	if err != nil {
		log.Println("Error parsing content read from stdin")
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
	//log.Printf("Response body\n%v\n", resp.Body)
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
	bodyStr, _ = json.Marshal(body)
	//log.Printf("bodyStr\n%v\n", string(bodyStr))
	actual.Body = bodyStr

	var headerStr json.RawMessage
	headerStr, _ = json.Marshal(headers)
	actual.Headers = headerStr

	return actual
}
func compareActualVersusExpected(actual actual, expect expect) bool {
	// compare the actual response against the expected response, and return a
	// boolean indicating whether the match is good or bad

	switch expect.ParseAs {
	case "":
		// if there's no parser defined, return false; this is a viable approach
		// for simply collecting info about an API response but not a test case
		return false
	case "regex":
		return true
	case "exact_match":
		return false
	case "partial_match":
		return false
	default:
		panic("expect.parse_as should be one of 'regex', 'exact_match', 'partial_match'")
	}
}

func main() {
	testCaseRequest := readTestJSON()

	testinfo, request, expect := populateRequest(testCaseRequest)
	body, headers, httpCode, latency := executeRequest(request)
	actual := populateResponse(body, headers, httpCode, latency)

	var passFail string
	if compareActualVersusExpected(actual, expect) {
		passFail = "pass"
	} else {
		passFail = "fail"
	}

	testresult := &testResult{
		PassFail:  passFail,
		Timestamp: time.Now().Local().Format(time.RFC3339),
		Request:   request,
		TestInfo:  testinfo,
		Expect:    expect,
		Actual:    actual,
	}

	testresultJSON, _ := json.MarshalIndent(testresult, "", "  ")
	fmt.Printf("%+v\n", string(testresultJSON))
}
