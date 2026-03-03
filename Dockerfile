FROM golang:alpine AS builder
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o Tiktok

FROM scratch
WORKDIR /app
COPY --from=builder /build/Tiktok /app
EXPOSE 8888
CMD ["./Tiktok"]
