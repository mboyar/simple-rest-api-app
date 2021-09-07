#build stage
FROM golang:bullseye AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod init app
RUN go mod download 
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM debian:bullseye
COPY --from=builder /go/bin/app /app
#RUN apt update && apt install -y systemd-container
ENTRYPOINT /app
LABEL Name=simplerestapiapp Version=0.0.1
EXPOSE 8080
