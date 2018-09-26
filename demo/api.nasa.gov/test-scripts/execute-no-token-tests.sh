#!/bin/bash
TESTCASES=`./test-scripts/select-only-no-token-tests.sh`
#APP=https://api.nasa.gov TESTCASE=`./test-scripts/select-only-no-token-tests.sh | xargs` MAX_CONCURRENT=3 go run ../../wilee/main.go
for f in $TESTCASES
do
    cat $f | APP=https://api.nasa.gov go run ../../wilee/main.go
done
