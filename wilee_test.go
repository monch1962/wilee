package main

import (
	"log"
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

func TestLoadValidJSON(t *testing.T) {
	testJSONfile := "demo/jsonplaceholder.typicode.com/test-cases/jsonplaceholder-test.json"
	fileHandle, err := os.Open(testJSONfile)
	if err != nil {
		t.Logf("test data file %v not found", testJSONfile)
		t.Fail()
	}
	_, err = readTestCaseJSON(fileHandle)
	if err != nil {
		t.Logf("Unable to parse test data file %v as JSON", testJSONfile)
		t.Fail()
	}
}

func TestLoadBrokenJSON(t *testing.T) {
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
