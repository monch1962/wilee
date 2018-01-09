#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
APP="https://jsonplaceholder.typicode.com" go run ../wilee.go < test-cases/jsonplaceholder-test.json | jq '.pass_fail'
