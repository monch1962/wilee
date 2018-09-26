#!/bin/bash
TAG=$1
rm test-cases/*.result.json 2> /dev/null
TEST_CASES=`ls test-cases/*.json`
for TC in $TEST_CASES
do
  TAGS=`cat $TC | jq -c '.test_info.tags'`
  if [[ $TAGS =~ $TAG ]]
  then
    echo "$TC"
  fi
done
