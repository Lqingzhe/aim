FROM golang:1.26.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download || echo "download failed, will use local cache if available"

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o ai-service ./app/aiservice

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/ai-service .
COPY app/aiservice/config.yaml .


EXPOSE 8893
CMD ["./ai-service"]