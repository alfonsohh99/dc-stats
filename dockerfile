# Dockerfile
FROM golang:latest
RUN mkdir /dc-stats
ADD . /dc-stats
WORKDIR /dc-stats
RUN go build -o bot .
CMD ["/dc-stats/bot"