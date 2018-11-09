#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

#RANDOM_HEADER_KEY=$(base64 -i /dev/urandom | fold -w 32 | head -1)
RANDOM_HEADER_KEY=$(dd if=/dev/random bs=256 count=1 2>/dev/null|od -An -tx1| tr -d ' \t\n')

RANDOM_HEADER_VALUE=$(base64 -i /dev/urandom | fold -w 4096 | head -1)
#echo $RANDOM_HEADER_KEY
#echo $RANDOM_HEADER_VALUE

echo $(cat) |
   jq --arg RHK "$RANDOM_HEADER_KEY" --arg RHV "$RANDOM_HEADER_VALUE" '.request.payload.headers += [{"header": ($RHK|tostring), "value": ($RHV|tostring)}]' |
   #jq '.expect.http_code=400' |
   jq '.test_info.tags -= ["huge_extra_header"]' |  
   jq '.test_info.tags += ["huge_extra_header"]'
