FROM golang:latest

WORKDIR /go/src/GoCrypt
COPY . .

RUN go get -d -v ./...
# Skip compiling files for now.
# RUN go install -v ./...


# Skip launching services for now.
#cmd ["app"]

# Allow API to be exposed.
EXPOSE 8000
