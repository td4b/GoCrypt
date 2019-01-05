FROM golang:latest
COPY . /usr/local/src/GoCrypt
EXPOSE 8000
RUN ["/bin/bash"]
