FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-meta-redirector .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/go-meta-redirector /usr/local/bin/go-meta-redirector
COPY repos.yaml /etc/go-meta-redirector/repos.yaml
EXPOSE 8080
ENTRYPOINT ["go-meta-redirector"]
