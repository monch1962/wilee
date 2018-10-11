#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# Extracting XML body for processing with xq can be problematic. 
echo $(cat) | 
    # Parsing out the XML using jq is pretty easy...
    jq '.expect.body' | 

    # xq doesn't like any initial <?xml version = ...?> so we have to get rid of everything from the first byte to ?>
    sed 's/^.*?>//' | 

    # Now we need to scrub off the trailing " character
    sed 's/.$//' | 

    # Finally remove any escape characters, as xq doesn't like them
    sed 's/\\//g' | 
    xq '.'