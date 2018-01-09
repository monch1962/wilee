#!/bin/bash
TEST_CASE="test-cases/jsonplaceholder-test.json"
PERF_TEST_ENV="http://www.loadtestenv.com"
TEST_WARMUP_DURATION_MINS=5
TEST_WARMUP_HIT_RATE=5
TEST_DURATION_MINS=20
TEST_INITIAL_HIT_RATE=300
TEST_FINAL_HIT_RATE=400

VERB=`cat $TEST_CASE | jq '.request.verb'`
#echo $VERB

cat $TEST_CASE | jq --arg testenv $PERF_TEST_ENV '.config.target |= $testenv' \
| jq --arg testwarmupduration $TEST_WARMUP_DURATION_MINS '.config.phases[0].duration |= $testwarmupduration' \
| jq --arg testwarmuphitrate $TEST_WARMUP_HIT_RATE '.config.phases[0].arrivalRate |= $testwarmuphitrate' \
| jq '.config.phases[0].name |= "Warm up period"' \
| jq --arg testduration $TEST_DURATION_MINS '.config.phases[1].duration |= $testduration' \
| jq --arg testinitialhitrate $TEST_INITIAL_HIT_RATE '.config.phases[1].arrivalRate |= $testinitialhitrate' \
| jq --arg testfinalhitrate $TEST_FINAL_HIT_RATE '.config.phases[1].rampTo |= $testfinalhitrate' \
| jq '.config.phases[1].name |= "Simple load test; no parameter substitution or chaining"' \
| jq '.scenarios[0].name |= "Auto-generated Artillery load test"' \
| jq '.scenarios[0].flow.url = .request.url' \
| jq '.scenarios[0].flow.body = .request.body' \
| jq 'del(.expect, .test_info, ._comment, .request)' > /tmp/tc

./json2yaml.py /tmp/tc | sed -e "s/'//g" | sed '/null/d'
