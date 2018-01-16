#!/bin/bash
POSTMAN_TESTS=./test-cases/postman-collections/*.json

for postmanfile in ${POSTMAN_TESTS}
do
  TEST_ID=`cat $postmanfile | jq -r '.info._postman_id'`
  TEST_DESCRIPTION=`cat $postmanfile | jq '.info.name'`
  #echo "TEST_ID: " $TEST_ID
  #echo "TEST_DESCRIPTION: " $TEST_DESCRIPTION
  POSTMAN_HOST=`cat $postmanfile | jq -r '.item[0].request.url.raw'`
  #echo $POSTMAN_HOST
  # extract the protocol
  proto="$(echo $POSTMAN_HOST | grep :// | sed -e's,^\(.*://\).*,\1,g')"
  #echo "PROTO: " $proto
  # remove the protocol
  url="$(echo ${POSTMAN_HOST/$proto/})"
  echo "URL: " $url

    # extract the path (if any)
  path="$(echo $url | grep / | cut -d/ -f2- | sed -es,^,/,)"
  #echo "PATH: " $path
  hostname="$(echo ${POSTMAN_HOST/$path/})"
  echo "HOSTNAME: " $hostname

  REQUEST_VERB=`cat $postmanfile | jq -r '.item[0].request.method'`
	cat $postmanfile \
    | jq --arg tc $postmanfile '._comment |= $tc' \
    | jq '.test_info.tags[0] |= "postman"' \
    | jq -r --arg testid $TEST_ID '.test_info.id = $testid' \
    | jq -r --arg targethost $hostname '.test_info.postman_host = $targethost' \
    | jq --arg requestverb $REQUEST_VERB '.request.verb = $requestverb' \
    | jq --arg url $path '.request.url = $url' \
    | jq 'del(.info,.item)'
    #| jq --arg testdescription $TEST_DESCRIPTION '.test_info.description |= $testdescription'

  #echo "TEST_ID:" $TEST_ID
  #cat $tc | jq '.info.name'
done
