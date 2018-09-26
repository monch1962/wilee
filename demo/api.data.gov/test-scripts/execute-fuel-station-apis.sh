#!/bin/bash

# Here we're going to take a test case, and mutate it 3 different ways:
# - first we'll add an API key to a header, then execute the test
# - then we'll add an API key to the GET request, and execute the test
# - then we'll leave out the API key, say we'll expect to get a HTTP 403 response, and execute the test

TESTCASES=`./test-scripts/select-tests-by-tag.sh api-key-required`
API_TOKEN=`cat test-cases/DATA_GOV_API_KEY`
for f in $TESTCASES
do
    # Append a header to the request called X-Api-Key and give it the API token as a value
    cat $f | jq --arg T "$API_TOKEN" '.request.payload.headers += [{"key": "X-Api-Key", "value": ($T)}]' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
done

for f in $TESTCASES
do
    # Grab the .request.url field out of each test case and append "?api_key=$API_TOKEN" to it before running it
    cat $f | jq --arg T "$API_TOKEN" '.request.url = (.request.url + "&api_key=" + $T)' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
done

for f in $TESTCASES
do
    # We're not passing an API key, so the return code should be a 403 - check this...
    cat $f | jq '.expect.http_code = 403' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.'
done

#for f in $TESTCASES
#do
    # We're not passing an API key, so the return code should be a 403 - check this...
#    cat $f | jq --arg T "$API_TOKEN" '.request.payload.body = "test"' | APP=https://developer.nrel.gov go run ../../wilee/main.go | jq '.request'
#done
