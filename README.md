# wilee (JSON test runner)

App designed to execute REST API functional test cases that have been encoded as JSON documents

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
