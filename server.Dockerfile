FROM golang:1.17.8 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go get gopkg.in/yaml.v2

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# multistage build to copy only binary and config
FROM scratch

COPY --from=builder /build/main /
COPY --from=builder /build/config/config.yaml /config/config.yaml

#EXPOSE 3333

ENTRYPOINT ["/main"]
