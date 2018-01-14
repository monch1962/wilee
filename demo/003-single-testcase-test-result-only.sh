#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
APP="http://localhost:51062" go run ../wilee.go < test-cases/jsonplaceholder-test.json | jq '.pass_fail'
