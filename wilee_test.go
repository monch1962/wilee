package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

/*func TestAssembleHTTPParameterString(t *testing.T) {
	var paramValues parameter
	log.Printf("%s\n", paramValues)
	paramValues = append(paramValues, map[string][string]{"key": ["value"]})
	log.Printf("%s\n", paramValues)

	//t.Fail()
}*/

func TestHelp(t *testing.T) {
	displayHelp()
	t.Log("displayHelp() is executing OK")
}

func TestLogResponseHeaders(t *testing.T) {
	var resp http.Response
	//resp.Headers.Set("abc", "123")
	logResponseHeaders(&resp)
	t.Log("logResponseHeaders() is executing OK")
}

func TestPopulateHTTPRequestHeaders(t *testing.T) {
	var req http.Request
	var headers []header
	os.Setenv("DEBUG", "1")
	req2 := populateHTTPRequestHeaders(&req, headers)
	//log.Printf("req2: %v\n", req2)
	if req2 != &req {
		t.Log("populateHTTPRequestHeaders() not working - adding no headers changes request")
		t.Fail()
	}
	var h header
	h.Header = "content-type"
	h.Value = "application/json"
	headers = []header{h}
	req2 = populateHTTPRequestHeaders(&req, headers)
	//t.Log("here2")
	//log.Printf("req2: %v\n", req2)
	if req2 != &req {
		t.Log("populateHTTPRequestHeaders() not working - adding no headers changes request")
		t.Fail()
	}
}

func TestPopulateRequest(t *testing.T) {
	var tc testCase
	tc.TestInfo.ID = "abc"
	tc.TestInfo.Description = "def"
	tc.TestInfo.Version = "1.99"
	tc.Request.Verb = "GET"
	testInfo, request, expect, _ := populateRequest(tc)
	log.Printf("testInfo: %v\n", testInfo)
	log.Printf("request: %v\n", request)
	log.Printf("expect: %v\n", expect)
	if testInfo.ID != "abc" {
		t.Log("populateRequest() not working - testInfo.ID not being populated")
		t.Fail()
	}
	if testInfo.Description != "def" {
		t.Log("populateRequest() not working - testInfo.Description not being populated")
		t.Fail()
	}
	if testInfo.Version != "1.99" {
		t.Log("populateRequest() not working - testInfo.Version not being populated")
		t.Fail()
	}
	if request.Verb != "GET" {
		t.Log("populateRequest() not working - request.Verb not being populated")
		t.Fail()
	}
}
func TestValidateIdenticalHTTPcodes(t *testing.T) {
	var expect expect
	var actual actual
	expect.HTTPCode = 200
	actual.HTTPCode = 200
	if !validateHTTPcodes(expect, actual) {
		t.Log("validateHTTPcodes not working - doesn't return true when codes are identical")
		t.Fail()
	}
}

func TestDebugOn(t *testing.T) {
	os.Setenv("DEBUG", "1")
	if !debug() {
		t.Log("debug() not working - doesn't return true when DEBUG='1'")
		t.Fail()
	}
}

func TestStringInArrayTrue(t *testing.T) {
	s := "abc"
	arr := []string{"34", "abc", "56"}
	if !stringInArray(s, arr) {
		t.Log("stringInArray() not working - doesn't return true when string is in array")
		t.Fail()
	}
}

func TestStringInArrayFalse(t *testing.T) {
	s := "abc"
	arr := []string{"34", "def", "56"}
	if stringInArray(s, arr) {
		t.Log("stringInArray() not working - doesn't return false when string isn't in array")
		t.Fail()
	}
}

func TestMaxConcurrency(t *testing.T) {
	os.Setenv("MAX_CONCURRENT", "")
	if maxConcurrency() != 1 {
		t.Log("maxConcurrency() not working - doesn't return 1 when MAX_CONCURRENT not set")
		t.Fail()
	}
	os.Setenv("MAX_CONCURRENT", "12")
	if maxConcurrency() != 12 {
		t.Log("maxConcurrency() not working - doesn't return 12 when MAX_CONCURRENT set to '12'")
		result := maxConcurrency()
		fmt.Printf("maxConcurrency: %v\n", result)
		t.Fail()
	}
}

/*func TestUnmarshalActualBody(t *testing.T) {
	var a actual
	log.Printf("a: %v\n", a)
	response, _ := unmarshalActualBody(a)
	log.Println("here3")
	if response != nil {
		t.Log("unmarshalActualBody() not working - doesn't return correct result")
		t.Fail()
	}
	raw := json.RawMessage(`{"foo":"bar"}`)
	var err error
	a.Body, err = json.Marshal(&raw)
	if err != nil {
		log.Printf("error marshalling body")
	}
	log.Printf("a: %v\n", a)
	response, _ := unmarshalActualBody(a)
	log.Println("here3")
	if response != "" {
		t.Log("unmarshalActualBody() not working - doesn't return correct result")
		t.Fail()
	}
}*/

func TestAssembleHTTPParamString(t *testing.T) {
	os.Setenv("DEBUG", "1")
	var parameters []parameter
	log.Printf("parameters: %v\n", parameters)
	if assembleHTTPParamString(parameters) != "" {
		t.Log("assembleHTTPParamString() not working - incorrect response when no parameters supplied")
		t.Fail()
	}
	var p1 parameter
	p1.Key = "abc"
	p1.Value = []string{"123"}
	log.Printf("p1: %v\n", p1)
	//parameters.append(parameters, parameter)
	parameters = []parameter{p1}
	if assembleHTTPParamString(parameters) != "?abc=123" {
		log.Printf("%s\n", assembleHTTPParamString(parameters))
		t.Log("assembleHTTPParamString() not working - incorrect response when parameter supplied")
		t.Fail()
	}
	var p2 parameter
	p2.Key = "def"
	p2.Value = []string{"456"}
	log.Printf("p2: %v\n", p2)
	parameters = []parameter{p1, p2}
	if assembleHTTPParamString(parameters) != "?abc=123&def=456" {
		log.Printf("%s\n", assembleHTTPParamString(parameters))
		t.Log("assembleHTTPParamString() not working - incorrect response when parameter supplied")
		t.Fail()
	}
}

func TestDebugOff(t *testing.T) {
	os.Setenv("DEBUG", "")
	if debug() {
		t.Log("debug() not working - doesn't return false when DEBUG=''")
	}
}
func TestValidateDifferentHTTPcodes(t *testing.T) {
	var expect expect
	var actual actual
	expect.HTTPCode = 200
	actual.HTTPCode = 300
	if validateHTTPcodes(expect, actual) {
		t.Log("validateHTTPcodes not working - doesn't return false when codes are different")
		t.Fail()
	}
}

func TestValidateMaxLatencyTrue(t *testing.T) {
	var expect expect
	var actual actual
	expect.MaxLatencyMS = 50
	actual.LatencyMS = 30
	if !validateMaxLatency(expect, actual) {
		t.Log("validateMaxLatency() not working - it doesn't return true with expect.MaxLatencyMS < actual.LatencyMS")
		t.Fail()
	}
}

func TestValidateMaxLatencyFalse(t *testing.T) {
	var expect expect
	var actual actual
	expect.MaxLatencyMS = 30
	actual.LatencyMS = 50
	if validateMaxLatency(expect, actual) {
		t.Log("validateMaxLatency() not working - it doesn't return false with expect.MaxLatencyMS > actual.LatencyMS")
		t.Fail()
	}
}

func TestLoadValidJSON(t *testing.T) {
	testJSONfile := "demo/jsonplaceholder.typicode.com/test-cases/jsonplaceholder-test.json"
	fileHandle, err := os.Open(testJSONfile)
	if err != nil {
		t.Logf("test data file %v not found", testJSONfile)
		t.Fail()
	}
	os.Setenv("DEBUG", "1")
	_, err = readTestCaseJSON(fileHandle)
	if err != nil {
		t.Logf("Unable to parse test data file %v as JSON", testJSONfile)
		t.Fail()
	}
}

func TestLoadBrokenJSON(t *testing.T) {
	os.Setenv("DEBUG", "1")
	brokenJSONfile := "demo/jsonplaceholder.typicode.com/test-cases/invalid/broken-json.json"
	fileHandle, err := os.Open(brokenJSONfile)
	if err != nil {
		t.Logf("broken test data file %v not found", brokenJSONfile)
		t.Fail()
	}
	_, err = readTestCaseJSON(fileHandle)
	if err == nil {
		t.Logf("Failed to detect broken JSON test file %v", brokenJSONfile)
		t.Fail()
	}
}

func TestLoadInvalidRequestJSON(t *testing.T) {
	brokenJSONfile := "demo/jsonplaceholder.typicode.com/test-cases/invalid/invalid-request-json.json"
	fileHandle, err := os.Open(brokenJSONfile)
	if err != nil {
		t.Logf("broken test data file %v not found", brokenJSONfile)
		t.Fail()
	}
	testcase, err := readTestCaseJSON(fileHandle)
	if err != nil {
		t.Logf("Unable to parse test data file %v as JSON", brokenJSONfile)
		t.Fail()
	}
	log.Printf("testcase: %v", testcase.Request)
	/*if testcase.Request != "" {
		t.Logf("Test case file %v should have invalid request")
		t.Fail()
	}*/
}

//func TestLoadInvalidJSON(t *testing.T) {
//	t.Fail()
//}
