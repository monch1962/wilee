#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

RANDOM_KEY=$(base64 -i /dev/urandom | fold -w 20 | head -1)
RANDOM_VALUE=$(base64 -i /dev/urandom | fold -w 20 | head -1)
echo $RANDOM_KEY
echo $RANDOM_VALUE

echo $(cat) |
   jq --arg RHK "$RANDOM_KEY" --arg RHV "$RANDOM_VALUE" '.request.body += {($RHK|tostring): ($RHV|tostring)}' |
   jq '.expect.http_code=400' |
   jq '.test_info.tags += ["random_extra_body_field"]'