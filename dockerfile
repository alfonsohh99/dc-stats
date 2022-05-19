# Dockerfile
FROM golang:latest
RUN mkdir /vc-stats
ADD . /vc-stats
WORKDIR /vc-stats
RUN go build -o bot .
CMD ["/vc-stats/bot"