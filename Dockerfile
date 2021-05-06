FROM golang:1.16.3-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git openssh

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server"]
