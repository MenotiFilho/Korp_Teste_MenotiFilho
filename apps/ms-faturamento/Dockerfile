FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/ms-faturamento ./cmd/ms-faturamento

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates wget

COPY --from=builder /bin/ms-faturamento /app/ms-faturamento

EXPOSE 8082

ENTRYPOINT ["/app/ms-faturamento"]
