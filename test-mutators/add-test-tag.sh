#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

TAG=$1

# First try to remove the tag in case it already exists (don't want to duplicate an existing tag, then add it)
echo $(cat) | 
    jq --arg TAG $1 '.test_info.tags -= [$TAG]' |
    jq --arg TAG $1 '.test_info.tags += [$TAG]'