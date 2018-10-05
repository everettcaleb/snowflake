FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/everettcaleb/snowflake/
RUN go get -v github.com/bronze1man/yaml2json
RUN go get -v github.com/gin-gonic/gin

COPY specs specs
COPY server.go server.go
COPY helpers.go helpers.go
COPY snowflake.go snowflake.go
RUN yaml2json > specs/spec.json < specs/spec.yaml
RUN go get .
RUN go build -o server .

FROM alpine:latest

ENV MACHINE_ID $HOST
ENV SNOWFLAKE_EPOCH 1514764800
ENV APP_BASE_PATH /v1/snowflake
ENV PORT 8080

COPY --from=builder /go/src/github.com/everettcaleb/snowflake/specs specs
COPY --from=builder /go/src/github.com/everettcaleb/snowflake/server server

EXPOSE 8080
CMD [ "./server" ]
