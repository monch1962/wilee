#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"


echo $(cat) |
   jq --arg RHK "$RANDOM_HEADER_KEY" --arg RHV "$RANDOM_HEADER_VALUE" '.request.headers[$RHK] = $RHV' |
   jq '.request.url=.request.url+"/-1"' |
   jq '.expect.http_code=400' |
   jq '.test_info.tags -= ["negative"]' |
   jq '.test_info.tags += ["negative"]' |
   jq '.test_info.tags += ["append_number_to_request_url"]'