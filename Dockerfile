FROM golang:alpine
RUN apk add --no-cache git
WORKDIR /go/src/github.com/everettcaleb/snowflake/
RUN go get -v github.com/bronze1man/yaml2json
RUN go get -v github.com/gorilla/mux

COPY specs specs
COPY server.go server.go
RUN cat specs/spec.yaml | yaml2json > specs/spec.json
RUN go build -o server .

FROM alpine

ENV MACHINE_ID $HOST
ENV SNOWFLAKE_EPOCH 1514764800
ENV APP_BASE_PATH /v1/snowflake
ENV PORT 8080

COPY --from=0 /go/src/github.com/everettcaleb/snowflake/specs specs
COPY --from=0 /go/src/github.com/everettcaleb/snowflake/server server

EXPOSE 8080
CMD [ "./server" ]
