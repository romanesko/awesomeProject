FROM golang:1.21.6-alpine3.19 AS builder

RUN apk update && apk upgrade && apk add --no-cache openssh=9.6_p1-r0

WORKDIR /go/src/app
COPY src .
RUN go get -u ./... && go build && go install

FROM scratch

COPY --from=builder /go/bin/server /

CMD ["/server"]
