FROM golang:1.20-alpine3.17 as builder
COPY . /src
WORKDIR /src
ENV GOPROXY "https://goproxy.cn"
RUN go build -o /build/to-read .

FROM alpine:3.17 as prod
COPY --from=builder /build/to-read /usr/bin/to-read
WORKDIR /app
ENTRYPOINT [ "to-read" ]
