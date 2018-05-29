FROM alpine:latest
# Minimal Docker container for hosting wilee
# (based on Alpine to make it as small as possible)
#
# To build:
# $ GOARCH=amd64 GOOS=linux go build wilee.go && docker build -t wilee:latest . && rm wilee
#
# To debug within built container:
# $ docker run -ti -e APP=https://jsonplaceholder.typicode.com -v `pwd`/test-data:/test-data -e TESTCASE=test-data/jsonplaceholder-test[0124].json --entrypoint /bin/ash wilee:latest
# # ./wilee
#
# Typical use case:
# $ docker run -d -e APP=https://jsonplaceholder.typicode.com -v `pwd`:/tests -e TESTCASE=tests/jsonplaceholder-test[0124].json wilee:latest
ENV UPDATED_AT 2018-01-03
ENV WILEE_HOME /opt/wilee
ENV TEST_CASES $WILEE_HOME/test_cases

RUN mkdir -p $WILEE_HOME
RUN mkdir -p $TEST_CASES

COPY ./wilee $WILEE_HOME/wilee

ENTRYPOINT ["/bin/ash", "$WILEE_HOME/wilee < $TEST_CASES/*"]
