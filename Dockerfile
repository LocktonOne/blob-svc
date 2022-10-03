FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/blob_api
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/blob_api /go/src/blob_api


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/blob_api /usr/local/bin/blob_api
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["blob_api"]
