package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
	//Content json.RawMessage `json:"content"`
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
	HTTPCode     int         `json:"http_code"`
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

// populateRequest takes the content of the test case and parses it into
// the JSON fragment that will eventually be returned from the test run
func populateRequest(testCaseRequest testRequest) (testInfo, request, expect) {
	testinfo := &testInfo{
		ID:          testCaseRequest.TestInfo.ID,
		Description: testCaseRequest.TestInfo.Description,
		Version:     testCaseRequest.TestInfo.Version,
		DateUpdated: testCaseRequest.TestInfo.DateUpdated,
		Author:      testCaseRequest.TestInfo.Author,
		//Content: testCaseRequest.TestInfo.Content,
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

// populateResponse takes info about the response and populates the "response" fragment of the JSON output
func populateResponse(body interface{}, headers interface{}, httpCode int, latency time.Duration) actual {
	var actualResponse actual
	actualResponse.HTTPCode = httpCode
	actualResponse.LatencyMS = int64(latency / time.Millisecond)

	var bodyStr json.RawMessage
	bodyStr, err := json.Marshal(body)
	if err != nil {
		panic("Unable to parse response body as JSON")
	}
	actualResponse.Body = bodyStr

	var headerStr json.RawMessage
	headerStr, err = json.Marshal(headers)
	if err != nil {
		panic("Unable to parse response headers as JSON")
	}
	actualResponse.Headers = headerStr

	return actualResponse
}

// compareActualVersusExpected compares the actual response against the
// expected response, and returns a boolean indicating whether the match was
// good or bad
func compareActualVersusExpected(actual actual, expect expect) bool {
	//log.Printf("expect.HTTPCode:%d\n", expect.HTTPCode)
	if expect.HTTPCode != 0 {
		if expect.HTTPCode != actual.HTTPCode {
			return false
		}
	}
	if expect.MaxLatencyMS != 0 {
		if expect.MaxLatencyMS < actual.LatencyMS {
			return false
		}
	}
	switch expect.ParseAs {
	case "":
		// if there's no parser defined, return false; this is a viable approach
		// for simply collecting info about an API response but not a test case
		return false
	case "regex":
		// we want to parse the actual content against regex patterns that are defined
		// in the "expected" part of the test case
		//log.Printf("expect.Body:%v\n", expect.Body)
		var b interface{}
		err := json.Unmarshal(actual.Body, &b)
		if err != nil {
			panic("Unable to parse actual.Body")
		}
		//log.Printf("actual.Body:%v\n", b)
		for k, expectRegex := range expect.Body.(map[string]interface{}) {
			log.Printf("expect[%s]->%v\n", k, expectRegex)
			actualValue := b.(map[string]interface{})[k]
			log.Printf("actual[%s]->%v\n", k, actualValue)
			log.Printf("actual.(type): %T\n", actualValue)

			r, _ := regexp.Compile(expectRegex.(string))
			switch actualValue.(type) {
			case int:
				if r.MatchString(string(actualValue.(int))) != true {
					return false
				}
			case float64:
				if r.MatchString(fmt.Sprintf("%f", actualValue.(float64))) != true {
					return false
				}
			case string:
				if r.MatchString(string(actualValue.(string))) != true {
					return false
				}
			}
		}

		return true
	case "exact_match":
		// we want the actual response to be an exact match to the "expected" part of
		// the test case. If there are extra fields in the actual response to what's
		// defined in "expected", then the test should fail
		var actualBodyStruct interface{}
		err := json.Unmarshal(actual.Body, &actualBodyStruct)
		if err != nil {
			panic("Unable to parse actual.Body")
		}
		for k, expectValue := range expect.Body.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := actualBodyStruct.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false
			}
		}
		for k, expectValue := range actualBodyStruct.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := expect.Body.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false
			}
		}
		return true
	case "partial_match":
		// we want the actual response fields to be an exact match to the "expected" fields
		// defined in the test case, but the "expected" fields may not contain all the fields in
		// the actual response
		//log.Printf("expect.Body:%v\n", expect.Body)
		var b interface{}
		err := json.Unmarshal(actual.Body, &b)
		if err != nil {
			panic("Unable to parse actual.Body")
		}
		for k, expectValue := range expect.Body.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := b.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false
			}
		}
		return true
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
