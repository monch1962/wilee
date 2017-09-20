#!/bin/bash
ALTERED_URL=/gefiltefish
cat test-data/jsonplaceholder-test.json | jq --arg url "$ALTERED_URL" '.request.url |= $url' | jq '.expect.body.body |= "Nothing good can come from this!"' | jq '.expect.max_latency_ms |= 1'
