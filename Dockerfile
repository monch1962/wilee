FROM alpine:latest
# Minimal Docker container for hosting jtrunner
# (based on Alpine to make it as small as possible)
#
# To build:
# $ GOARCH=amd64 GOOS=linux go build jtrunner.go && docker build -t jtrunner:latest . && rm jtrunner
#
# To debug within built container:
# $ docker run -ti -e APP=https://jsonplaceholder.typicode.com -v `pwd`/test-data:/test-data -e TESTCASE=test-data/jsonplaceholder-test[0124].json --entrypoint /bin/ash jtrunner:latest
# # ./jtrunner
#
# Typical use case:
# $ docker run -d -e APP=https://jsonplaceholder.typicode.com -v `pwd`:/tests -e TESTCASE=tests/jsonplaceholder-test[0124].json jtrunner:latest
ENV UPDATED_AT 2017-05-24
COPY ./jtrunner jtrunner
ENTRYPOINT ["./jtrunner"]
