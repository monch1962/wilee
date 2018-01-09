#!/bin/bash
TEST_CASE="test-cases/jsonplaceholder-test.json"
PERF_TEST_ENV="http://www.loadtestenv.com"
TEST_DURATION_MINS=20
TEST_INITIAL_HIT_RATE=300
TEST_FINAL_HIT_RATE=400

VERB=`cat $TEST_CASE | jq '.request.verb'`
#echo $VERB

cat $TEST_CASE | jq --arg testenv $PERF_TEST_ENV '.config.target |= $testenv' \
| jq --arg testduration $TEST_DURATION_MINS '.config.phases[0].duration |= $testduration' \
| jq --arg testinitialhitrate $TEST_INITIAL_HIT_RATE '.config.phases[0].arrivalRate |= $testinitialhitrate' \
| jq --arg testfinalhitrate $TEST_FINAL_HIT_RATE '.config.phases[0].rampTo |= $testfinalhitrate' \
| jq '.config.phases[0].name |= "Simple load test; no parameter substitution or chaining"' \
| jq '.scenarios[0].name |= "Auto-generated Artillery load test"' \
| jq '.scenarios[0].flow.url = .request.url' \
| jq '.scenarios[0].flow.body = .request.body' \
| jq 'del(.expect, .test_info, ._comment, .request)' > /tmp/tc

./json2yaml.py /tmp/tc | sed -e "s/'//g" | sed '/null/d'
