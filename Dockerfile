FROM golang:latest as builder

WORKDIR /usr/src/
COPY main.go go.mod ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /usr/src/url-shortener .
CMD ["./url-shortener"]
