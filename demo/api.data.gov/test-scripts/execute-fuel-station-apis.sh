#!/bin/bash
TESTCASES=`./test-scripts/select-tests-by-tag.sh api-key-required`
API_TOKEN=`cat test-cases/DATA_GOV_API_KEY`
for f in $TESTCASES
do
    # Append a header to the request called X-Api-Key and give it the API token as a value
    cat $f | jq --arg T "$API_TOKEN" '.request.payload.headers += [{"X-Api-Key": ($T)}]' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
done

#for f in $TESTCASES
#do
    # Grab the .request.url field out of each test case and append "?api_key=$API_TOKEN" to it before running it
#    cat $f | jq --arg T "$API_TOKEN" '.request.url = (.request.url + "&api_key=" + $T)' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
#done

#for f in $TESTCASES
#do
    # We're not passing an API key, so the return code should be a 403 - check this...
#    cat $f | jq '.expect.http_code = 403' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
#done

for f in $TESTCASES
do
    # We're not passing an API key, so the return code should be a 403 - check this...
    cat $f | jq --arg T "$API_TOKEN" '.request.payload.body = "test"' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.request'
done
