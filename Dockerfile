FROM golang:1.12-alpine

WORKDIR /go/src/app
COPY . .

RUN apk add --update ca-certificates
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -mod vendor -o werds

CMD ./werds
