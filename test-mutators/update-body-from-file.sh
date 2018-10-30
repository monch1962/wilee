#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"


echo $(cat) | 
    jq --slurpfile BODY_FROM_FILE $1 '.request.body = $BODY_FROM_FILE[0]'