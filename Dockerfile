FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
# Skip compiling files for now.
# RUN go install -v ./...

cmd ["go run cryptapi.go"]


# Skip launching services for now.
#cmd ["app"]
