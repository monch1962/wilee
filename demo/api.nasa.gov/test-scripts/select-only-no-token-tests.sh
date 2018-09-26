#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
TEST_CASES=`ls test-cases/*.json`
#echo $TEST_CASES
#echo -n '{'
for TC in $TEST_CASES
do
  TAGS=`cat $TC | jq -c '.test_info.tags'`
  #echo $TAGS
  if [[ $TAGS =~ "no_token" ]]
  then
    echo "$TC"
  fi
done
#echo "}"
