#!/usr/bin/env python
import sys
import json
import yaml # need to 'pip install pyyaml' for this to work; 'brew install libyaml && sudo python -m easy_install pyyaml' on Mac

print (yaml.dump(yaml.load(json.dumps(json.loads(open(sys.argv[1]).read()))), default_flow_style=False))
