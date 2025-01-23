FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/uploader main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/bin/uploader /app/uploader
COPY config.yaml /app/config.yaml
EXPOSE 8081
ENTRYPOINT ["/app/uploader"]
