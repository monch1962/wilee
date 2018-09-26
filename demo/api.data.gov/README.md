Here we're hitting APIs on the https://api.data.gov site

In this case, we've got a single test case that hits https://developer.nrel.gov/api/alt-fuel-stations/v1.json?limit=1
This URL requires an API key to be provided, but there's 2 different ways in which we can provide it:
- via a request header using the key X-Api-Key
- via "&api_key=..." appended to the GET request

The test script test-scripts/execute-fuel-station-apis.sh takes the basic test case in test-cases/fuel-stations.json and mutates it in 3 different ways:
- first we append a header to the basic test JSON to contain the API key, then execute the modified test case
- next we append "&api_key=..." to the GET statement, then execute the modified test case
- finally we leave out the API key altogether, and change the expected http response code to 403, then execute the modified test case