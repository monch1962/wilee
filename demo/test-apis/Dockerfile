# Test JSON API server for wilee
#
# To build
# $ docker build -t wilee-test-apis .
#
# To start
# $ docker run -p 34567:3000 -d wilee-test-apis

FROM node:alpine
WORKDIR /app
COPY apis/apis.json .
RUN npm install -g json-server
ENTRYPOINT ["json-server", "--watch", "apis.json"]
EXPOSE 3000
