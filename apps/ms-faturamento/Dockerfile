FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/ms-faturamento ./cmd/ms-faturamento

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/ms-faturamento /app/ms-faturamento

EXPOSE 8082

ENTRYPOINT ["/app/ms-faturamento"]