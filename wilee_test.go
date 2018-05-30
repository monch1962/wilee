package main

import (
	"log"
	"os"
	"testing"
)

func TestLoadValidJSON(t *testing.T) {
	testJSONfile := "demo/test-cases/jsonplaceholder-test.json"
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
	brokenJSONfile := "demo/test-cases/invalid/broken-json.json"
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
	brokenJSONfile := "demo/test-cases/invalid/invalid-request-json.json"
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
