#!/bin/bash
TESTCASES=`./test-scripts/select-tests-by-tag.sh neo`
API_TOKEN=`cat test-cases/NASA_API_KEY`
for f in $TESTCASES
do
    # Grab the .request.url field out of each test case and append "?api_key=$API_TOKEN" to it before running it
    cat $f | jq --arg T "$API_TOKEN" '.request.url = (.request.url + $T)' | APP=https://api.nasa.gov go run ../../wilee/main.go | jq '.'
done
