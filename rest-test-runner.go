package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Header struct {
	Header string `json:"header"`
	Value  string `json:"value"`
}

type TestInfo struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Version     string `json:"version"`
	DateUpdated string `json:"date_uploaded"`
	Author      string `json:"author"`
}

type Payload struct {
	Headers []Header `json:"headers"`
	Body    string   `json:"body"`
}

type Request struct {
	Verb    string  `json:"verb"`
	Url     string  `json:"url"`
	Payload Payload `json:"payload"`
}

type Expect struct {
	ParseAs      string      `json:"parse_as"`
	HttpCode     int64       `json:"http_code"`
	MaxLatencyMS int64       `json:"max_latency_ms"`
	Headers      []Header    `json:"headers"`
	Body         interface{} `json:"body"`
}

type Actual struct {
	HttpCode  int64    `json:"http_code"`
	LatencyMS int64    `json:"latency_ms"`
	Headers   []Header `json:"headers"`
	Body      string   `json:"body"`
}

type TestResult struct {
	PassFail  string   `json:"pass_fail"`
	Timestamp string   `json:"timestamp"`
	TestInfo  TestInfo `json:"test_info"`
	Request   Request  `json:"request"`
	Expect    Expect   `json:"expect"`
	Actual    Actual   `json:"actual"`
}

type TestRequest struct {
	TestInfo TestInfo `json:"test_info"`
	Request  Request  `json:"request"`
	Expect   Expect   `json:"expect"`
}

type TestInput struct {
	testinfo TestInfo `json:"test_info"`
}

func readTestJson() TestRequest {
	// Read JSON input from stdin and return as a formatted Go struct
	j, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("Error reading content from stdin")
		panic(err)
	}
	var pj TestRequest
	err = json.Unmarshal(j, &pj)
	if err != nil {
		log.Println("Error parsing content read from stdin")
		log.Printf("%v\n", string(j))
		panic(err)
	}
	//log.Printf("Go struct: %v", pj)
	formattedInput, _ := json.MarshalIndent(pj, "", "  ")
	log.Printf("Input JSON: %+v\n", string(formattedInput))

	return pj
}

func main() {
	testCaseRequest := readTestJson()

	log.Printf("testCase:\n%+v\n", testCaseRequest)

	testinfo := &TestInfo{
		Id:          testCaseRequest.TestInfo.Id,
		Description: testCaseRequest.TestInfo.Description,
		Version:     testCaseRequest.TestInfo.Version,
		DateUpdated: testCaseRequest.TestInfo.DateUpdated,
		Author:      testCaseRequest.TestInfo.Author,
	}

	request := &Request{
		Verb: testCaseRequest.Request.Verb,
		Url:  testCaseRequest.Request.Url,
	}

	expect := &Expect{
		ParseAs:      testCaseRequest.Expect.ParseAs,
		HttpCode:     testCaseRequest.Expect.HttpCode,
		MaxLatencyMS: testCaseRequest.Expect.MaxLatencyMS,
		Headers:      testCaseRequest.Expect.Headers,
		Body:         testCaseRequest.Expect.Body,
	}

	testresult := &TestResult{
		PassFail:  "pass",
		Timestamp: time.Now().Local().Format(time.RFC3339),
		Request:   *request,
		TestInfo:  *testinfo,
		Expect:    *expect,
	}

	testresultJSON, _ := json.MarshalIndent(testresult, "", "  ")
	fmt.Printf("%+v\n", string(testresultJSON))
}
