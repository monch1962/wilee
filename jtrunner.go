package main

// Test runner for functional tests defined as JSON documents
// Expects the following environment variables to be defined:
// APP = target server for requests e.g. "https://SERVER:PORT"
// STUB_ENGINE = type of stub engine e.g. "stubby", "montebank"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type header struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type testInfo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	DateUpdated string   `json:"date_uploaded"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
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
	PassFail       string   `json:"pass_fail"`
	PassFailReason string   `json:"pass_fail_reason"`
	Timestamp      string   `json:"timestamp"`
	TestInfo       testInfo `json:"test_info"`
	Request        request  `json:"request"`
	Expect         expect   `json:"expect"`
	Actual         actual   `json:"actual"`
}

type testCase struct {
	TestInfo testInfo `json:"test_info"`
	//TODO: StubConfig
	Request request `json:"request"`
	Expect  expect  `json:"expect"`
}

// readTestCaseJSON reads a JSON testcase from an io.Reader and returns it as a formatted Go
// struct
func readTestCaseJSON(input io.Reader) (testCase, error) {
	j, err := ioutil.ReadAll(input)
	var ti testCase
	if err != nil {
		return ti, errors.New("Error reading JSON test case content")
	}
	err = json.Unmarshal(j, &ti)
	if err != nil {
		return ti, errors.New("Error parsing content as JSON")
	}
	return ti, nil
}

// if the test case requires stub configuration, do so
func configureStubEngine(testCase) {
	stubEngine := os.Getenv("STUB_ENGINE")
	log.Printf("Sending configuration to stub engine %v\n", stubEngine)
	switch strings.ToUpper(stubEngine) {
	case "MONTEBANK":
		log.Println("Stub engine is Montebank")
	case "STUBBY":
		log.Println("Stub engine is Stubby")
	default:
		log.Printf("Don't know what to do for stub engine '%v'\n", stubEngine)
	}
}

// populateRequest takes the content of the test case and parses it into
// the JSON fragment that will eventually be returned from the test run
func populateRequest(tc testCase) (testInfo, request, expect, error) {
	testinfo := &testInfo{
		ID:          tc.TestInfo.ID,
		Description: tc.TestInfo.Description,
		Version:     tc.TestInfo.Version,
		DateUpdated: tc.TestInfo.DateUpdated,
		Author:      tc.TestInfo.Author,
		//Content: tc.TestInfo.Content,
	}

	request := &request{
		Verb: tc.Request.Verb,
		URL:  os.Getenv("APP") + tc.Request.URL,
	}

	expect := &expect{
		ParseAs:      tc.Expect.ParseAs,
		HTTPCode:     tc.Expect.HTTPCode,
		MaxLatencyMS: tc.Expect.MaxLatencyMS,
		Headers:      tc.Expect.Headers,
		Body:         tc.Expect.Body,
	}
	return *testinfo, *request, *expect, nil
}

// executeRequest executes the JSON request defined in the test case, and captures & returns
// the response body, response headers, HTTP status and latency
func executeRequest(request request) (interface{}, interface{}, int, time.Duration, error) {
	if request.Verb != "GET" && request.Verb != "POST" && request.Verb != "PUT" && request.Verb != "DELETE" && request.Verb != "HEAD" && request.Verb != "PATCH" {
		return nil, nil, 0, 0, errors.New("request.verb must be one of GET, POST, PUT, DELETE, HEAD, PATCH")
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest(request.Verb, request.URL, nil)
	if err != nil {
		return nil, nil, 0, 0, errors.New("Unable to parse HTTP request")
		//log.Fatalln(err)
	}
	startTime := time.Now()
	resp, err := httpClient.Do(req)
	endTime := time.Since(startTime)
	if err != nil {
		return nil, nil, 0, 0, errors.New("Unable to execute HTTP request")
	}
	defer resp.Body.Close()
	responseDecoder := json.NewDecoder(resp.Body)
	var v interface{} // Not sure what the response will look like, so just implement an interface
	err = responseDecoder.Decode(&v)
	if err != nil {
		return nil, nil, 0, 0, errors.New("Unable to parse HTTP response body as JSON")
	}
	//log.Printf("v\n%v\n", v)
	//log.Println(resp.Header)
	//for i := range resp.Header {
	//log.Printf("%v->%v\n", i, resp.Header[i][0])
	//}
	headers := resp.Header
	httpCode := resp.StatusCode
	latency := endTime
	return v, headers, httpCode, latency, nil
}

// populateResponse takes info about the response and populates the "response" fragment of the JSON output
func populateResponse(body interface{}, headers interface{}, httpCode int, latency time.Duration) (actual, error) {
	var actualResponse actual
	actualResponse.HTTPCode = httpCode
	actualResponse.LatencyMS = int64(latency / time.Millisecond)

	var bodyStr json.RawMessage
	bodyStr, err := json.Marshal(body)
	if err != nil {
		return actualResponse, errors.New("Unable to parse response body as JSON")
	}
	actualResponse.Body = bodyStr

	var headerStr json.RawMessage
	headerStr, err = json.Marshal(headers)
	if err != nil {
		return actualResponse, errors.New("Unable to parse response headers as JSON")
	}
	actualResponse.Headers = headerStr

	return actualResponse, nil
}

// compareActualVersusExpected compares the actual response against the
// expected response, and returns a boolean indicating whether the match was
// good or bad
func compareActualVersusExpected(actual actual, expect expect) (bool, string, error) {
	//log.Printf("expect.HTTPCode:%d\n", expect.HTTPCode)
	if expect.HTTPCode != 0 {
		if expect.HTTPCode != actual.HTTPCode {
			return false, "actual.HTTPCode doesn't match", nil
		}
	}
	if expect.MaxLatencyMS != 0 {
		if expect.MaxLatencyMS < actual.LatencyMS {
			return false, "actual.latency_ms > expect.max_latency_ms", nil
		}
	}
	switch expect.ParseAs {
	case "":
		// if there's no parser defined, return false; this is a viable approach
		// for simply collecting info about an API response but not a test case
		return false, "expect.parse_as not defined", nil
	case "regex":
		// we want to parse the actual content against regex patterns that are defined
		// in the "expected" part of the test case
		//log.Printf("expect.Body:%v\n", expect.Body)
		var b interface{}
		err := json.Unmarshal(actual.Body, &b)
		if err != nil {
			return false, "", errors.New("Unable to parse actual.Body")
		}
		//log.Printf("actual.Body:%v\n", b)
		for k, expectRegex := range expect.Body.(map[string]interface{}) {
			log.Printf("expect[%s]->%v\n", k, expectRegex)
			actualValue := b.(map[string]interface{})[k]
			log.Printf("actual[%s]->%v\n", k, actualValue)
			log.Printf("actual.(type): %T\n", actualValue)

			r, err := regexp.Compile(expectRegex.(string))
			if err != nil {
				log.Fatalf("Error compiling regex for expect.%s", k)
			}
			switch actualValue.(type) {
			case int:
				if r.MatchString(string(actualValue.(int))) != true {
					return false, "expect.* doesn't match actual.*", nil
				}
			case float64:
				if r.MatchString(fmt.Sprintf("%f", actualValue.(float64))) != true {
					return false, "", nil
				}
			case string:
				if r.MatchString(string(actualValue.(string))) != true {
					return false, "", nil
				}
			}
		}

		return true, "", nil
	case "exact_match":
		// we want the actual response to be an exact match to the "expected" part of
		// the test case. If there are extra fields in the actual response to what's
		// defined in "expected", then the test should fail
		var actualBodyStruct interface{}
		err := json.Unmarshal(actual.Body, &actualBodyStruct)
		if err != nil {
			//panic("Unable to parse actual.Body")
		}
		for k, expectValue := range expect.Body.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := actualBodyStruct.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false, "expectValue != actualValue", nil
			}
		}
		for k, expectValue := range actualBodyStruct.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := expect.Body.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false, "expectValue != actualValue", nil
			}
		}
		return true, "", nil
	case "partial_match":
		// we want the actual response fields to be an exact match to the "expected" fields
		// defined in the test case, but the "expected" fields may not contain all the fields in
		// the actual response
		//log.Printf("expect.Body:%v\n", expect.Body)
		var b interface{}
		err := json.Unmarshal(actual.Body, &b)
		if err != nil {
			//panic("Unable to parse actual.Body")
		}
		for k, expectValue := range expect.Body.(map[string]interface{}) {
			//log.Printf("expect[%s]->%v\n", k, expectValue)
			actualValue := b.(map[string]interface{})[k]
			//log.Printf("actual[%s]->%v\n", k, actualValue)
			if expectValue != actualValue {
				//log.Printf("expectValue != actualValue: %v -> %v", expectValue, actualValue)
				return false, "expectValue != actualValue", nil
			}
		}
		return true, "", nil
	default:
		return false, "", errors.New("expect.parse_as should be one of 'regex', 'exact_match', 'partial_match'")
	}
}

func executeTestCase(testFile *os.File, resultsFile *os.File) {
	// read the JSON test case from file
	tc, err := readTestCaseJSON(testFile)
	if err != nil {
		log.Println("Unable to read test case JSON input")
		os.Exit(1)
	}

	// if there's a stub to be configured, do so
	if os.Getenv("STUB_ENGINE") != "" {
		log.Println("STUB_ENGINE configured")
		configureStubEngine(tc)
	}

	// populate the the "request" content that will eventually be sent to stdout
	testinfo, request, expect, err := populateRequest(tc)
	if err != nil {
		log.Println("Unable to parse test request info out of test case")
		os.Exit(1)
	}

	// execute the request, and capture the response body, headers, http status and latency
	body, headers, httpCode, latency, err := executeRequest(request)
	if err != nil {
		log.Println("Unable to execute HTTP request")
		os.Exit(1)
	}

	// populate the "response" content that will eventually be sent to stdout
	actual, err := populateResponse(body, headers, httpCode, latency)
	if err != nil {
		log.Println("Unable to populate HTTP response JSON")
		os.Exit(1)
	}

	// check whether the "response" matches what was expected, which defines whether
	// the test run passed or failed
	var passFail = "fail"
	matchSuccess, passFailReason, err := compareActualVersusExpected(actual, expect)
	if err != nil {
		log.Println("Unable to compare actual response vs. expected response")
		os.Exit(1)
	}
	if matchSuccess {
		passFail = "pass"
	}

	// construct the output JSON
	testresult := &testResult{
		PassFail:       passFail,
		PassFailReason: passFailReason,
		Timestamp:      time.Now().Local().Format(time.RFC3339),
		Request:        request,
		TestInfo:       testinfo,
		Expect:         expect,
		Actual:         actual,
	}

	// make the output JSON look pretty
	testresultJSON, err := json.MarshalIndent(testresult, "", "  ")
	if err != nil {
		//panic("Unable to display output as JSON")
	}

	// send the output JSON to stdout
	//fmt.Printf("%+v\n", string(testresultJSON))
	resultsFile.WriteString(string(testresultJSON))
}

func executeTestCaseInWaitGroup(testFile *os.File, resultsFile *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	executeTestCase(testFile, resultsFile)
}

func main() {
	if os.Getenv("TESTCASE") != "" {
		testcases := os.Getenv("TESTCASE")
		testCaseFilesGlob, _ := filepath.Glob(testcases)
		//log.Printf("Test cases to execute: %v\n", testCaseFilesGlob)
		numTestCases := len(testCaseFilesGlob)
		log.Printf("# test cases to execute = %d\n", numTestCases)

		var wg sync.WaitGroup
		//wg.Add(numTestCases)
		for _, testCaseFile := range testCaseFilesGlob {
			log.Printf("Running test case from file: %v\n", testCaseFile)
			fTestCase, err := os.Open(testCaseFile)
			if err != nil {
				log.Printf("Error opening test case file %v\n", testCaseFile)
			}
			testResultsFile := testCaseFile + ".result.json"
			fTestResults, err := os.Create(testResultsFile)
			if err != nil {
				log.Printf("Error opening test case results file %v\n", testResultsFile)
			}
			wg.Add(1)
			go executeTestCaseInWaitGroup(fTestCase, fTestResults, &wg)
		}
		wg.Wait()
		//time.Sleep(5 * time.Second)
	} else {
		executeTestCase(os.Stdin, os.Stdout)
	}
}
