# test-mutators

Scripts in this directory are designed to mutate an existing test case in some way. They are written to expect a test case to be supplied through stdin and to emit the mutated test case to stdout

Note that most of these scripts will assume jq is installed, and that they're running under a bash-like shell

Try `cat TESTCASE.JSON | add-tags "mexico"`
