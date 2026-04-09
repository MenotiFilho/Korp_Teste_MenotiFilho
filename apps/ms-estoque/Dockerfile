FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/ms-estoque ./cmd/ms-estoque

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/ms-estoque /app/ms-estoque

EXPOSE 8081

ENTRYPOINT ["/app/ms-estoque"]
