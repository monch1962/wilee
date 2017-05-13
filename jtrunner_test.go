package main

import (
	"os"
	"testing"
)

func TestLoadValidJSON(t *testing.T) {
	testJSONfile := "test-data/jsonplaceholder-test.json"
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

func TestLoadInvalidJSON(t *testing.T) {
	t.Fail()
}
