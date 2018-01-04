# wilee

App designed to execute REST API functional test cases that have been encoded as JSON documents

## Quick start
Assuming you've got Go installed:
```
$ git clone https://github.com/monch1962/wilee
$ cd wilee
$ go build wilee.go
$ APP="https://jsonplaceholder.typicode.com" wilee < test-data/jsonplaceholder-test.json
```

## Why do integration testing?

Look, I can't stand integration testing (commonly referred to as SIT, for System Integration Testing). It's fragile, slow, massively expensive, and tends to chew up loads of time and require loads of people. Lots of organisations are moving away from performing any SIT, and most others are probably wishing they could.

However, for a lot of organisations, SIT is still seen as essential. There's still a risk of changes somewhere breaking something somewhere else, and that risk has to be managed somehow.

In most workplaces I've been at, SIT is generally conducted as follows:
* you've got a bunch of developers who build your solution in some combination of (programming language/toolsets/frameworks). Let's call these PTFs
* employ a bunch of "expert" testers. These people will set up their test environments using some different combination of PTFs and then write your SIT tests. This combination of PTFs might be similar to those used by your developers, but there will be test-specific PTFs in there that your developers won't be using. Alternately, the SIT PTFs might be vastly different from those used by your developers. Either way, you'll have more "stuff" to deal with
* the end result is that you'll have a bunch of SIT test artefacts - that's "test environment" plus "code" plus "data" plus "test results" plus maybe some documentation - that you'll have to maintain over the life of the project

### And...?
Here's some harsh lessons I've learned with respect to the above:
* _good_ testers are hard to find and tend to be expensive; you really want them to be as productive as possible, as quickly as possible
* most of the value of _good_ testers comes in their ability to identify _what_ needs to be tested, work out _how_ to test it & explain _why_ exactly that set of things needs to be tested. Their ability to actually create test cases is somewhere further down the list
* installing an automation test toolset should be as fast, simple and idiot-proof as possible, on any platform from a tester's laptop to a CI/CD pipeline. Very few tools meet those criteria
* people talk about test-driven development (TDD) as though it's easily achievable; it _should be possible to create new automation tests within seconds or minutes_, not hours
  * adding new test cases should be a no-brainer, near-zero-time task that can be performed before development starts
  * with practice, a scrum team should be able to create new automation tests _during_ sprint planning (WTF?)
* test cases should be highly _mutable by design_; test cases are almost never static, and need to evolve over the lifetime of a project
* test cases should be _atomic_ (e.g. single file per test case, with no external dependencies)
* people tend to focus on test case _creation_, but the tough part is what comes after that: test lifecycle management
* test execution should be _extremely_ fast
* test execution should normally be performed as part of a CI/CD pipeline, automatically triggered via e.g. a git branch update
* test execution should be highly _auditable_ and support any audit framework
* creating and maintaining _test suites_ should be a trivial exercise
* avoid test frameworks as far as possible, as they tie you into a specific way of working
* testing tools should support current _and future_ best practice processes

## What's wrong with other SIT tools in this space?

This is always going to be a highly opinionated list, but here goes...
* high level of expertise required to create automation test cases
* some come with a large set of dependencies
* focused primarily on test *creation*, rather than lifecycle management, productivity, agility, maintaining flexibility, auditability, future proofing, CI/CD friendliness, ...
* tend to lock you into to specific processes, toolsets, frameworks, programming languages, etc.
* can be a challenge to interface with test management tools such as JIRA

## What does wilee bring to the table?

wilee sets out to be the *smallest possible viable framework* for API integration testing. It seeks to isolate you from all PTF-related issues, which means you have almost nothing "new" to maintain besides the stuff your developers are using. The entire wilee "framework" is one file - the wilee executable you built above

wilee seeks to utilise the fantastic command-line tools that have been developed since Unix first appeared in the 1960s. Tools like 'bash', 'cat', 'jq' in particular; these tools are free, they have just about every conceivable bug shaken out, and your developers probably already know how to use them.

wilee makes it easy to leverage the brilliance that is Docker containers, along with whatever container orchestration framework you happen to be using (Docker Compose, Kubernetes, etc.). Docker makes loads of stuff easier, including SIT. wilee won't force you to use Docker, but Docker makes life a lot easier.

wilee tries to make it simple to integrate with whatever test management tools you might be using - that's JIRA, ALM and the like - using bash and wilee itself.

## Design decisions

wilee is deliberately opinionated, which has driven the following design decisions:
* test cases are written in JSON, and test execution results are presented in JSON
* wilee is written in Go
* command line interface only
* minimal feature sets

Now let's look at the implications of these decisions...

### Why JSON?
JSON is:
* language and toolset agnostic - "everything understands JSON"
* easy for humans to write and reason about - "have I get test coverage for this?"
* test results are easy for computers to parser
* JSON actually makes a great DSL for writing test cases & capturing test results
* tests become atomic (1 JSON per test case, 1 JSON per test case execution)
* easy to mutate existing test cases to create new tests (e.g. generative testing, fuzzy testing)
* easy to evolve existing test cases to deal with new data, changes in field content, adding/deleting fields, etc.
* opens up new options for managing test lifecycle
* easy to extend
* provides a path for repurposing functional automation tests as performance tests

### Why Go?
Go:
* compiles down to a single, small executable file on any platform, with no additional runtime
  * makes it easy to setup and reset "test runner" environments
  * makes it very easy to deploy into containers, which is where I prefer to use for test environments
  * makes it easy to integrate with CI/CD pipelines
* creates executables that have very small RAM and disk footprints
* is very fast to execute (e.g. no startup delay as with JVM languages)
* offers native support for multithreading via goroutines
* cross compilation is built into the language

### Why command line interface only?
* Every OS supports execution from the command line
* Easy to scale from running a single test to running 1000s of tests
* Easy to leverage different test case & results storage options

### Why a minimal feature set?
wilee attempts to follow the Unix philosophy as documented at https://en.Wikipedia.org/wiki/Unix_philosophy
* “Write programs that do one thing and do it well”
* “Expect the output of every program to become the input to another, as yet unknown, program”
* “Write programs to work together”
* “Write programs to handle text streams, because that is a universal interface”
* “Use shell scripts to increase leverage and portability”
* “Developers should design programs so that they do not print unnecessary output”
* “Developers should design for the future by making their protocols extensible”


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Download and Compile

```
$ git clone https://github.com/monch1962/wilee
$ cd wilee
$ go build wilee.go
```

You should then see a new wilee executable file. It'll be named 'wilee.exe' if you're running Windows, or 'wilee' if you're running pretty much anything else.

You may want to built wilee for a different system to that you're currently using - for example, you might be sitting in front of a Mac, but you want to build wilee to run within a Linux Docker container. In that case, change the last command to:
```
$ OS=linux ARCH=amd64 go build wilee.go
```
and you'll get a file named 'wilee' that'll run on your Linux container. It *WON'T* run on your Mac, because the binary won't be compatible.

### Quick test

To quickly test your new wilee executable, try running on your non-Windows box
```
$ APP="https://jsonplaceholder.typicode.com" wilee < test-data/jsonplaceholder-test.json
```

Assuming you've got Internet access, this will run a test against the https://jsonplaceholder.typicode.com site, and return a blob of JSON. Feel free to scroll through the output to see what's there...

### Prerequisites

You'll probably want to define a set of environment variables before executing tests using wilee.

An important one is to set the target app server that you wish to test

```
APP="https://server:port"
```
wilee can either read a single test case from stdin and write results to stdout (where they can be processed by e.g. jq), or you can supply a set of test cases to be executed. In the latter case, all test cases will be executed simultaneously
```
TESTCASE=test-data/jsonplaceholder-test[0124].json
```
### Installing

A step by step series of examples that tell you have to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Running the tests

Explain how to run the automated tests for this system

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With


## Contributing

## Versioning

## Authors

* **David Mitchell** - *Initial work* - (https://github.com/monch1962)

See also the list of [contributors](https://github.com/monch1962/wilee/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone who's code was used
* Inspiration
* etc
