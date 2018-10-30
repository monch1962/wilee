#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# Let's update some data in the database
UPDATE_DATA=$(sqlite3 testdata.db "update users set name='George Steele' where id=13;" ".exit")

# Now we'll extract that data and use it to update a test case
EXTRACT_DATA=$(sqlite3 testdata.db "select * from users;" ".exit")
ID=$(cut -d'|' -f1 <<<"$EXTRACT_DATA")
NAME=$(cut -d'|' -f2 <<<"$EXTRACT_DATA")

cat test.json | 
  jq --arg NAME "$NAME" '.expect.body.name=$NAME' |
  jq --arg ID "$ID" '.expect.body.id=$ID'

# Let's make another database update
UPDATE_DATA=$(sqlite3 testdata.db "update users set name='Fred Blassie' where id=13;" ".exit")

# Now we'll extract that data and use it to update the test case with different data
EXTRACT_DATA=$(sqlite3 testdata.db "select * from users;" ".exit")
ID=$(cut -d'|' -f1 <<<"$EXTRACT_DATA")
NAME=$(cut -d'|' -f2 <<<"$EXTRACT_DATA")

cat test.json | 
  jq --arg NAME "$NAME" '.expect.body.name=$NAME' |
  jq --arg ID "$ID" '.expect.body.id=$ID'

