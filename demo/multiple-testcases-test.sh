#!/bin/bash
rm test-data/*.result.json
time APP="https://jsonplaceholder.typicode.com" TESTCASE="test-data/jsonplaceholder-test[0124]*.json" go run ../jtrunner.go
for f in test-data/*.result.json 
do
	cat $f | jq '{result: .pass_fail, verb: .request.verb, url: .request.url }'
done

