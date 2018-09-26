#!/bin/bash
rm test-cases/*.result.json 2> /dev/null
TEST_CASES=`ls test-cases/*.json`
#echo $TEST_CASES
for TC in $TEST_CASES
do
  TAGS=`cat $TC | jq -c '.test_info.tags'`
  #echo $TAGS
  if [[ $TAGS =~ "token_required" ]]
  then
    echo -n "$TC "
  fi
done
echo
