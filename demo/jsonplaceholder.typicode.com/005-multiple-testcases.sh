#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
APP="http://localhost:51062" TESTCASES="test-cases/jsonplaceholder-test[012]*.json" go run ../wilee.go
for f in test-cases/*.result.json
do
	cat $f | jq '{result: .pass_fail, pass_fail_reason: .pass_fail_reason, verb: .request.verb, url: .request.url }' | jq --arg result "$f" '.test_result |= $result'
done
