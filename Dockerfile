FROM golang:1.21.6-alpine3.19

RUN apk update && apk upgrade && apk add --no-cache openssh=9.6_p1-r0

WORKDIR /go/src/app
COPY src .
RUN go get -u ./... && go build && go install

CMD ["/go/bin/server"]
