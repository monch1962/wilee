#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

RANDOM_HEADER_KEY=$(base64 -i /dev/urandom | fold -w 4096 | head -1)
RANDOM_HEADER_VALUE=$(base64 -i /dev/urandom | fold -w 4096 | head -1)
#echo $RANDOM_HEADER_KEY
#echo $RANDOM_HEADER_VALUE

echo $(cat) |
   jq --arg RHK "$RANDOM_HEADER_KEY" --arg RHV "$RANDOM_HEADER_VALUE" '.request.headers[$RHK] = $RHV' |
   #jq '.expect.http_code=400' |
   jq '.test_info.tags -= ["huge_extra_header"]' |  
   jq '.test_info.tags += ["huge_extra_header"]'