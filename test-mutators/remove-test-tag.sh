#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

TAG=$1

echo $(cat) | 
    jq --arg TAG $1 '.test_info.tags -= [$TAG]'