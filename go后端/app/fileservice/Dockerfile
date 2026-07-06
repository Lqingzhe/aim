FROM golang:1.26.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download || echo "download failed, will use local cache if available"

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o file-service ./app/fileservice

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/file-service .
COPY ./app/fileservice/config.yaml .


EXPOSE 8892
CMD ["./file-service"]