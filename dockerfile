# Dockerfile
FROM golang:latest
RUN mkdir /dc-stats
ADD . /dc-stats
WORKDIR /dc-stats
RUN go build -o /usr/local/bin/app .
CMD ["app"]