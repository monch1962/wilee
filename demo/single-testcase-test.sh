#!/bin/bash
time APP="https://jsonplaceholder.typicode.com" go run ../wilee.go < test-data/jsonplaceholder-test.json | jq
