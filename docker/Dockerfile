FROM golang:1.26.1 AS builder

ARG SERVICE

ENV TZ=Asia/Shanghai \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /app/output
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/output/${SERVICE} && \
    go build -o /app/output/${SERVICE}/tiktok-${SERVICE} ./cmd/${SERVICE}/

FROM alpine:3.21

ARG SERVICE

ENV TZ=Asia/Shanghai \
    SERVICE=${SERVICE}

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/output/${SERVICE}/tiktok-${SERVICE} /app/tiktok-${SERVICE}

CMD ["sh", "-c", "/app/tiktok-${SERVICE}"]
