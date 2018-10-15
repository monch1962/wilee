#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

echo $(cat) |
   jq '.expect.parse_as="partial_match"' |
   jq 'del(.expect.body)'|
   jq 'del(.expect.headers)'