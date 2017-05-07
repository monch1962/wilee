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
	HttpCode  int      `json:"http_code"`
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
	var tr TestRequest
	err = json.Unmarshal(j, &tr)
	if err != nil {
		log.Println("Error parsing content read from stdin")
		log.Printf("%v\n", string(j))
		panic(err)
	}

	return tr
}

func populateRequest(testCaseRequest TestRequest) (TestInfo, Request, Expect) {
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
	return *testinfo, *request, *expect
}

func executeRequest(request Request) (interface{}, interface{}, int, time.Duration) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(request.Verb, request.Url, nil)
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
	log.Printf("Response body\n%v\n", resp.Body)
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

func populateResponse(body interface{}, headers interface{}, httpCode int, latency time.Duration) Actual {
	var actual Actual
	actual.HttpCode = httpCode
	actual.LatencyMS = int64(latency / time.Millisecond)

	bodyStr, _ := json.Marshal(body)
	log.Printf("bodyStr\n%v\n", string(bodyStr))
	actual.Body = string(bodyStr)
	//b1, _ := json.MarshalIndent(bodyStr, "", "  ")
	//actual.Body = string(b1)

	//headerStr, _ := json.Marshal(headers)
	//actual.Headers = string(headerStr)

	return actual
}

func main() {
	testCaseRequest := readTestJson()

	log.Printf("testCase:\n%+v\n", testCaseRequest)

	testinfo, request, expect := populateRequest(testCaseRequest)

	body, headers, httpCode, latency := executeRequest(request)
	log.Printf("Response body\n%v\n", body)
	log.Printf("Response headers\n%v\n", headers)

	actual := populateResponse(body, headers, httpCode, latency)

	testresult := &TestResult{
		PassFail:  "pass",
		Timestamp: time.Now().Local().Format(time.RFC3339),
		Request:   request,
		TestInfo:  testinfo,
		Expect:    expect,
		Actual:    actual,
	}

	testresultJSON, _ := json.MarshalIndent(testresult, "", "  ")
	fmt.Printf("%+v\n", string(testresultJSON))
}
