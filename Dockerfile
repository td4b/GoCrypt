FROM golang:latest

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o main .

CMD ["/app/main"]

# Allow API to be exposed.
EXPOSE 8000
