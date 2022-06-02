# syntax=docker/dockerfile:1
FROM golang:1.17-alpine AS build_base

# Add Maintainer Info
LABEL maintainer="Sam Wang <sam.wang.0723@gmail.com>"

RUN apk add --no-cache git
RUN apk update && apk add ca-certificates && apk add tzdata

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v -tags all_tests -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v /third_party/)

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /stock-crawler cmd/main.go

# Start fresh from a smaller image
FROM scratch

#RUN apk add ca-certificates
#RUN apk add --no-cache tzdata

WORKDIR /

COPY --from=build_base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build_base /app/configs/config.dev.yaml /configs/config.dev.yaml
COPY --from=build_base /app/configs/stock_ids.json /configs/stock_ids.json
COPY --from=build_base /stock-crawler /stock-crawler

# This container exposes ports to the outside world
EXPOSE 8086 80 443 22 9092 6379

#USER nonroot:nonroot
#RUN chmod +x /jarvis-api
ENV TZ=Asia/Taipei

# Run the binary program
ENTRYPOINT [ "/stock-crawler" ]

