FROM golang:1.24-alpine AS builder

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

COPY . /src
WORKDIR /src

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/website-service ./cmd/website-service

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /src/bin /app
COPY --from=builder /src/configs /app/configs

WORKDIR /app

EXPOSE 8082
EXPOSE 9082

CMD ["./website-service", "-conf", "./configs/config.yaml"]
