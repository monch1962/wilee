#!/bin/bash
time APP="https://jsonplaceholder.typicode.com" go run ../jtrunner.go < test-data/jsonplaceholder-test.json | jq
