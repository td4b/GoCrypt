FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
# Skip compiling files for now.
# RUN go install -v ./...

RUN go build cryptapi.go

cmd ["./cryptapi.go"]


# Skip launching services for now.
#cmd ["app"]
