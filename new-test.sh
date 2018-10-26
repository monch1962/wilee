#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

VERB=$(echo $1|tr '[:lower:]' '[:upper:]')
URL=$2
HTTP_CODE=$3
TIMESTAMP=$(date +%Y-%m-%dT%H:%M:%S%z)
if [[ "$OSTYPE" == "darwin"* ]]; then
    USER=$(logname)
elif [[ "$OSTYPE" == "win"* ]]; then
    USER=%USERNAME%
fi

echo "{}" |
    jq --arg VERB $VERB --arg URL $URL '.test_info.id=($VERB)+" "+($URL)' |
    jq '.test_info.description = "wilee test case template"' |
    jq '.test_info.version = "0.01"' |
    jq --arg TIMESTAMP $TIMESTAMP '.test_info.date_uploaded=($TIMESTAMP)' |
    jq --arg USER $USER '.test_info.author=($USER)' |
    jq '.test_info._comment= "**GENERATED TEMPLATE - NEED TO UPDATE THIS TEST CASE AS IT WILL INITIALLY FAIL**"' |
    jq '.test_info.tags=[ "template_only"]' |
    jq '.request._comment=""' |
    jq --arg VERB $VERB '.request.verb=($VERB)' |
    jq --arg URL $URL '.request.url=($URL)' |
    jq '.request.headers=[]' |
    jq '.request.parameters=[]' |
    jq '.request.body=""' |
    jq --arg HTTP_CODE $HTTP_CODE '.expect.http_code=($HTTP_CODE|tonumber)' |
 #   jq '.expect.parse_as="exact_match"' |
    jq '.expect.max_latency_ms=0' 
