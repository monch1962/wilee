package main

// Test runner for functional tests defined as JSON documents
// Expects the following environment variables to be defined:
// APP = target server for requests e.g. "https://SERVER:PORT"
// TESTCASES (optional) = regex for a set of test case JSON files

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/monch1962/okerlund"
	"github.com/nsf/jsondiff"
	"github.com/remeh/sizedwaitgroup"
	"github.com/xeipuuv/gojsonschema"
)

type header struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type parameter struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type testInfo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	DateUpdated string   `json:"date_uploaded"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

type payload struct {
	Headers    []header    `json:"headers"`
	Body       interface{} `json:"body"`
	Parameters []parameter `json:"parameters"`
}

type request struct {
	Verb    string  `json:"verb"`
	URL     string  `json:"url"`
	Payload payload `json:"payload"`
}

type expect struct {
	ParseAs      string          `json:"parse_as"`
	HTTPCode     int             `json:"http_code"`
	MaxLatencyMS int64           `json:"max_latency_ms"`
	Headers      json.RawMessage `json:"headers"`
	Body         interface{}     `json:"body"`
}

type actual struct {
	HTTPCode  int             `json:"http_code"`
	LatencyMS int64           `json:"latency_ms"`
	Headers   json.RawMessage `json:"headers"`
	//Body      json.RawMessage `json:"body"`
	Body interface{} `json:"body"`
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
	Request  request  `json:"request"`
	Expect   expect   `json:"expect"`
}

type testCaseAwsAPIGatewayEvent struct {
	TestCase string `json:"queryStringParameters"`
}

func debug() bool {
	if os.Getenv("DEBUG") != "" {
		return true
	}
	return false
}

// readTestCaseJSON reads a JSON testcase from an io.Reader and returns it as a formatted Go
// struct
func readTestCaseJSON(input io.Reader) (testCase, error) {
	j, err := ioutil.ReadAll(input)
	var ti testCase
	if err != nil {
		log.Printf("Error reading JSON test case content")
		return ti, errors.New("Error reading JSON test case content")
	}
	err = json.Unmarshal(j, &ti)
	if err != nil {
		log.Printf("Error parsing content as JSON")
		return ti, errors.New("Error parsing content as JSON")
	}
	if debug() {
		log.Printf("testcase: %v\n", ti)
		log.Printf("testcase.request.payload.body: %v\n", ti.Request.Payload.Body)
		log.Printf("testcase.request.payload.headers: %v\n", ti.Request.Payload.Headers)
		log.Printf("testcase.request.payload.parameters: %v\n", ti.Request.Payload.Parameters)
	}
	return ti, nil
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
	}

	requestPayload := payload{
		Headers:    tc.Request.Payload.Headers,
		Body:       tc.Request.Payload.Body,
		Parameters: tc.Request.Payload.Parameters,
	}

	request := &request{
		Verb:    tc.Request.Verb,
		URL:     os.Getenv("APP") + tc.Request.URL,
		Payload: requestPayload,
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

func stringInArray(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func assembleHTTPParamString(parameters []parameter) string {
	var httpParamString strings.Builder
	for _, param := range parameters {
		if debug() {
			log.Printf("param: %s\n", param)
		}

		if httpParamString.String() == "" {
			httpParamString.WriteString("?")
		} else {
			httpParamString.WriteString("&")
		}
		httpParamString.WriteString(param.Key)
		httpParamString.WriteString("=")
		httpParamString.WriteString(html.UnescapeString(string(param.Value[0])))
	}
	if debug() {
		log.Printf("httpParamString: %s\n", httpParamString.String())
	}
	return httpParamString.String()
}

func populateHTTPRequestHeaders(req *http.Request, headers []header) *http.Request {
	for _, h := range headers {
		if debug() {
			log.Printf("request header: %v\n", h)
		}
		k := h.Header
		v := h.Value
		if debug() {
			log.Printf("request header key: %v\n", k)
			log.Printf("request header value: %v\n", v)
		}
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set(k, v)
	}
	return req
}

func logResponseHeaders(response *http.Response) {
	for i := range response.Header {
		log.Printf("%v->%v\n", i, response.Header[i][0])
	}
}

// executeRequest executes the JSON request defined in the test case, and captures & returns
// the response body, response headers, HTTP status and latency
func executeRequest(request request) (interface{}, interface{}, int, time.Duration, error) {
	validVerbs := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH"}
	if !stringInArray(request.Verb, validVerbs) {
		return nil, nil, 0, 0, errors.New("request.verb must be one of GET, POST, PUT, DELETE, HEAD, PATCH")
	}
	httpClient := &http.Client{}
	if debug() {
		log.Printf("req.Payload.Parameters: %v\n", request.Payload.Parameters)
	}
	var httpParamString string
	httpParamString = assembleHTTPParamString(request.Payload.Parameters)
	unescapedURL := request.URL + httpParamString
	if debug() {
		log.Printf("unescapedURL: %s\n", unescapedURL)
	}
	req, err := http.NewRequest(request.Verb, unescapedURL, nil)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, 0, 0, errors.New("Unable to parse HTTP request")
		//log.Fatalln(err)
	}

	req = populateHTTPRequestHeaders(req, request.Payload.Headers)
	if debug() {
		log.Printf("req.Payload.Body: %s\n", request.Payload.Body)
	}
	if request.Payload.Body != nil && !reflect.ValueOf(request.Payload.Body).IsNil() {

		body, _ := json.Marshal(request.Payload.Body)
		if debug() {
			log.Printf("body: %s\n", body)
		}
		bodyReader := bytes.NewReader(body)
		bodyReadCloser := ioutil.NopCloser(bodyReader)
		req.Body = bodyReadCloser
		if debug() {
			log.Printf("httpRequest: %v\n", req)
		}
	}
	startTime := time.Now()
	resp, err := httpClient.Do(req)
	endTime := time.Since(startTime)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, 0, 0, errors.New("Unable to execute HTTP request")
	}
	defer resp.Body.Close()
	responseDecoder := json.NewDecoder(resp.Body)
	var v interface{} // Not sure what the response will look like, so just implement an interface
	err = responseDecoder.Decode(&v)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, 0, 0, errors.New("Unable to parse HTTP response body as JSON")
	}
	if debug() {
		log.Printf("v\n%v\n", v)
		log.Println(resp.Header)

		logResponseHeaders(resp)
	}
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

// JSONCompare compares 2 JSON strings: if expect is equivalent to actual, or expect is a subset of actual, then return true
func JSONCompare(actual []byte, expect []byte) jsondiff.Difference {
	opts := jsondiff.DefaultConsoleOptions()
	opts.PrintTypes = true
	jsonDifference, _ := jsondiff.Compare(actual, expect, &opts)
	if debug() {
		log.Printf("JSONCompare difference: %s\n", jsonDifference)
	}
	return jsonDifference
}

func compareJSONSchema(expect expect, actual actual) bool {
	if expect.Body != nil {
		expectLoader := gojsonschema.NewGoLoader(expect)
		actualLoader := gojsonschema.NewGoLoader(actual)
		result, err := gojsonschema.Validate(expectLoader, actualLoader)
		if err != nil {
			log.Println("Error running JSON schema validation")
			panic(err.Error())
		}
		if !result.Valid() {
			fmt.Fprintln(os.Stderr, "JSON schema validation of response failed")
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
			return false
		}
	}
	return true
}

func validateHTTPcodes(expect expect, actual actual) bool {
	if expect.HTTPCode != 0 && expect.HTTPCode != actual.HTTPCode {
		return false
	}
	return true
}

func validateMaxLatency(expect expect, actual actual) bool {
	if expect.MaxLatencyMS != 0 && expect.MaxLatencyMS < actual.LatencyMS {
		return false
	}
	return true
}

func unmarshalActualBody(actual actual) (interface{}, error) {
	var actualBodyStruct interface{}
	err := json.Unmarshal(actual.Body.(json.RawMessage), &actualBodyStruct)
	if err != nil {
		return nil, errors.New("Unable to parse actual.Body")
	}
	if debug() {
		log.Printf("actual.Body:%s\n\n", actual.Body)
		log.Printf("actualBodyStruct: %v\n", actualBodyStruct)
	}
	return actualBodyStruct, nil
}

func compareRegex(expect expect, actual actual) (bool, string, error) {
	actualBodyStruct, err := unmarshalActualBody(actual)
	if err != nil {
		return false, "", errors.New("Unable to parse actual.Body")
	}
	if debug() {
		log.Printf("actual.Body:%s\n\n", actual.Body)
		log.Printf("actualBodyStruct: %v\n", actualBodyStruct)
	}
	if expect.Body != nil {
		for k, expectRegex := range expect.Body.(map[string]interface{}) {
			if debug() {
				log.Printf("expect[%s]->%v\n", k, expectRegex)
			}

			actualValue := actualBodyStruct.(map[string]interface{})[k]
			if debug() {
				log.Printf("actual[%s]->%v\n", k, actualValue)
				log.Printf("actual.(type): %T\n", actualValue)
			}

			r, err := regexp.Compile(expectRegex.(string))
			if err != nil {
				log.Fatalf("Error compiling regex for expect.%s", k)
			}

			var expectNotEqualActualMsg = "expect.* doesn't match actual.*"
			switch actualValue.(type) {
			case int:
				if !r.MatchString(string(actualValue.(int))) {
					return false, expectNotEqualActualMsg, nil
				}
			case float64:
				if !r.MatchString(fmt.Sprintf("%f", actualValue.(float64))) {
					return false, expectNotEqualActualMsg, nil
				}
			case string:
				if !r.MatchString(string(actualValue.(string))) {
					return false, expectNotEqualActualMsg, nil
				}
			}
		}
	}
	return true, "", nil
}

func compareJSON(expect expect, actual actual, comparisonType string) (bool, string, error) {
	if debug() {
		log.Printf("expect: %s\n", expect.Body)
		log.Printf("actual: %s\n", actual.Body)
	}
	if expect.Body != nil {
		expectJSON, _ := json.Marshal(expect.Body)
		difference := JSONCompare(actual.Body.(json.RawMessage), expectJSON)
		switch comparisonType {
		case "partial_match":
			if difference != jsondiff.FullMatch {
				return false, "expect.body is not a subset of actual.body", nil
			}
		case "exact_match":
			if difference != jsondiff.SupersetMatch {
				return false, "expect.body is not a subset of actual.body", nil
			}
		default:
			return false, "invalid comparison type - should be 'exact_match' or 'partial_match'", nil
		}
	}
	return true, "", nil
}

// compareActualVersusExpected compares the actual response against the
// expected response, and returns a boolean indicating whether the match was
// good or bad
func compareActualVersusExpected(actual actual, expect expect) (bool, string, error) {
	if !validateHTTPcodes(expect, actual) {
		errText := fmt.Sprintf("actual.HTTPCode doesn't match expect.HTTPCode. Expected %d, got %d", expect.HTTPCode, actual.HTTPCode)
		return false, errText, nil
	}
	if !validateMaxLatency(expect, actual) {
		errText := fmt.Sprintf("actual.latency_ms (%d) > expect.max_latency_ms (%d)", actual.LatencyMS, expect.MaxLatencyMS)
		return false, errText, nil
	}

	switch expect.ParseAs {
	case "":
		// if there's no parser defined, return false; this is a viable approach
		// for simply collecting info about an API response but not a test case
		return false, "expect.parse_as not defined", nil
	case "json_schema":
		if compareJSONSchema(expect, actual) {
			return true, "", nil
		}
		return false, "JSON schema validation of response failed", nil
	case "regex":
		// we want to parse the actual content against regex patterns that are defined
		// in the "expected" part of the test case
		//log.Printf("expect.Body:%v\n", expect.Body)
		return compareRegex(expect, actual)
	case "exact_match":
		// we want the actual response fields to be an exact match to the "expected" fields
		// defined in the test case, but the "expected" fields may not contain all the fields in
		// the actual response
		return compareJSON(expect, actual, "exact_match")
	case "partial_match":
		// we want the actual response fields to be an exact match to the "expected" fields
		// defined in the test case, but the "expected" fields may not contain all the fields in
		// the actual response
		return compareJSON(expect, actual, "partial_match")
	default:
		return false, "", errors.New("expect.parse_as should be one of 'regex', 'exact_match', 'partial_match', 'json_schema'")
	}
}

func executeTestCase(testFile *os.File, resultsFile *os.File) {
	// read the JSON test case from file
	tc, err := readTestCaseJSON(testFile)
	if err != nil {
		log.Printf("%s\n", err)
		log.Println("Unable to read test case JSON input")
		os.Exit(1)
	}

	testresultJSON, err := executeTestCaseJSON(tc)
	// send the output JSON to resultsFile
	resultsFile.WriteString(string(testresultJSON))
}

func executeTestCaseJSON(tc testCase) (testresultJSON string, err error) {
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
	testresultBytes, err := json.MarshalIndent(testresult, "", "  ")
	testresultJSON = string(testresultBytes)
	if err != nil {
		//panic("Unable to display output as JSON")
	}
	return testresultJSON, nil
}

func executeTestCaseInWaitGroup(testFile *os.File, resultsFile *os.File, wg *sizedwaitgroup.SizedWaitGroup) {
	defer wg.Done()
	executeTestCase(testFile, resultsFile)
}

//HandleRequest is the designated handler for Lambda
func HandleRequest(reqEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestStr, _ := json.Marshal(reqEvent.Body)
	log.Printf("request: %s\n", requestStr)
	tc, err := readTestCaseJSON(strings.NewReader(reqEvent.Body))
	if err != nil {
		log.Println("Unable to read test case JSON input")
		os.Exit(1)
	}

	testresultJSON, err := executeTestCaseJSON(tc)

	return events.APIGatewayProxyResponse{Body: string(testresultJSON), StatusCode: 200}, nil
}

func displayHelp() {
	fmt.Println()
	fmt.Println("wilee expects to see an environment variable APP pointing to the http/s server to be tested.")
	fmt.Println()
	fmt.Println("wilee reads a test case from stdin, and prints test case execution results to stdout")
	fmt.Println("This means that, if you simply run")
	fmt.Println("  $ wilee")
	fmt.Println("it will sit there forever, waiting for input")
	fmt.Println()
	fmt.Println("This probably isn't what you want...")
	fmt.Println()
	fmt.Println("Instead `cat TESTCASE.json | APP=http://localhost:8000 wilee` is a valid way to run a test against a server running at http://localhost:8000")
	fmt.Println()
	fmt.Println("Another option is to set the environment variable TESTCASES to point a set of test cases using wildcards, and use wilee to execute all of those test cases.")
	fmt.Println("For example,")
	fmt.Println("  $ APP=http://localhost:8000 TESTCASES=tests/test*.json wilee")
	fmt.Println("will run all test cases defined in tests/test*.json simultaneously")
	fmt.Println()
	fmt.Println("If you want to limit the number of test cases that are run concurrently, you can use the MAX_CONCURRENCY environment variable to do so.")
	fmt.Println("For example,")
	fmt.Println("  $ APP=http://localhost:8000 TESTCASES=tests/test*.json MAX_CONCURRENCY=3 wilee")
	fmt.Println("will run all test cases defined in tests/test*.json, but no more than 3 will run concurrently.")
}

func maxConcurrency() int {
	maxConcurrent, err := strconv.Atoi(os.Getenv("MAX_CONCURRENT"))
	if err != nil {
		return 1
	}
	return maxConcurrent
}

func executeRequestedTestcases() {
	// filenames for test cases to run are contained in the env var TESTCASES
	// which can contain regexes
	testcases := os.Getenv("TESTCASES")
	testCaseFilesGlob, _ := filepath.Glob(testcases)

	maxConcurrent := maxConcurrency()
	log.Printf("Max concurrency = %d\n", maxConcurrent)
	numTestCases := len(testCaseFilesGlob)
	log.Printf("# test cases to execute = %d\n", numTestCases)

	// we're going to run all these test cases in parallel in separate
	// goroutines...
	// ... but we also want to wait till all of them have finished before
	// exiting so we define a WaitGroup

	// Using sizedwaitgroup instead of sync.WaitGroup gives me an easy way to limit concurrency
	wg := sizedwaitgroup.New(maxConcurrent)

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

		// add this new test case to the wait group
		wg.Add()
		go executeTestCaseInWaitGroup(fTestCase, fTestResults, &wg)
	}

	// wait here till all test cases have finished executing
	wg.Wait()
}

func main() {
	switch okerlund.IsLambdaEnv() {
	case true:
		lambda.Start(HandleRequest)
	case false:
		//log.Printf("Not running in Lambda env\n")

		var helpPtr = flag.Bool("help", false, "Display help")
		flag.Parse()
		if *helpPtr {
			displayHelp()
			os.Exit(0)
		}
		if os.Getenv("TESTCASES") != "" {
			executeRequestedTestcases()
		} else {
			// no testcase files supplied in env var TESTCASES
			// read a single test case from stdin, and write test results to stdout
			executeTestCase(os.Stdin, os.Stdout)
		}
	}
}
