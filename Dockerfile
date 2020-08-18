FROM golang:latest as builder

WORKDIR /usr/src/
COPY main.go go.mod ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo


FROM heroku/heroku:latest

ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku
COPY --from=builder /usr/src/url-shortener /app/
CMD /app/url-shortener
