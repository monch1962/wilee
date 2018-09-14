#!/bin/bash
#RANDOM_CHARS="~!@#$%^&*()_+|}{:?,./;'[]\'"
#RANDOM=$$$(date +%s)
#RANDOM_CHAR=${RANDOM_CHARS[$RANDOM % ${#RANDOM_CHARS[@]}]}
RANDOM_CHAR=$(cat /dev/random | head -c 1)
echo $RANDOM_CHAR
echo $(cat) |
   jq --arg R "$RANDOM_CHAR" '.request.headers = [$R]' |
   jq '.expect.http_code=422' |
   jq '.test_info.tags += ["negative"]'|
   jq 'del(.expect.body)'