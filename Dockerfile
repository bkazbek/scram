FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/scram_tcp/
COPY . .

RUN go get -d -v

RUN export CGO_ENABLED=0 && go build -o /go/bin/scram_tcp

FROM alpine

COPY --from=builder /go/bin/scram_tcp /go/bin/scram_tcp

ENTRYPOINT ["/go/bin/scram_tcp"]