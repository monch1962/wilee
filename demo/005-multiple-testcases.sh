#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
time APP="https://jsonplaceholder.typicode.com" TESTCASE="test-cases/jsonplaceholder-test[0124]*.json" go run ../wilee.go
for f in test-cases/*.result.json
do
	cat $f | jq '{result: .pass_fail, verb: .request.verb, url: .request.url }' | jq --arg result "$f" '.test_result |= $result'
done
