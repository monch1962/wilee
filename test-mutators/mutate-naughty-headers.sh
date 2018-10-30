#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

LINES=$(cat $DIR/naughty-bodies/blns.txt | wc -l)
R_LINE=$((1+$RANDOM % $LINES))
randomnaughtystring=$(sed -n "${R_LINE}p" $DIR/naughty-bodies/blns.txt)

echo $(cat) |
   jq --arg RNS "$randomnaughtystring" '.request.headers = [{$RNS}]' |
   jq '.expect.http_code=422' |
   jq '.test_info.tags -= ["negative"]'|
   jq '.test_info.tags += ["negative"]'|
   jq 'del(.expect.body)'