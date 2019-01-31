FROM golang AS builder
ENV GO111MODULE on
WORKDIR /go/src/github.com/everettcaleb/snowflake/

COPY go.mod go.mod
COPY go.sum go.sum
RUN go get

COPY config.go config.go
COPY server.go server.go
COPY machine-id.go machine-id.go
COPY responses.go responses.go
COPY snowflake.go snowflake.go
RUN go build -o server .

FROM alpine:latest

ENV PORT 8080

COPY --from=builder /go/src/github.com/everettcaleb/snowflake/specs specs
COPY --from=builder /go/src/github.com/everettcaleb/snowflake/server server

EXPOSE 8080
CMD [ "./server" ]
