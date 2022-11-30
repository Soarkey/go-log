FROM golang:alpine AS builder
LABEL stage=gobuilder
WORKDIR /build

# RUN apk update --no-cache && apk add --no-cache tzdata
ENV GO111MODULE=on
ENV CGO_ENABLED 0
ENV GOPROXY=https://goproxy.cn,direct

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/main ./log-transfer/main.go

FROM alpine

# RUN apk update --no-cache && apk add --no-cache ca-certificates
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]