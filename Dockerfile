FROM golang:latest as builder

WORKDIR /usr/src/
COPY main.go go.mod ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo


FROM scratch

ENV PORT 6666
COPY --from=builder /usr/src/url-shortener /app/
ENTRYPOINT ["/app/url-shortener"]
