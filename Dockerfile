FROM golang:alpine as builder

RUN adduser --uid 2000 --disabled-password user

RUN apk add --no-cache git curl bash

COPY main.go /go/src/proxy/
COPY go.mod /go/src/proxy/
WORKDIR /go/src/proxy

RUN mkdir /out
RUN go build -o /out/proxy

EXPOSE 8080
USER 2000
ENV PORT=8080
CMD /out/proxy
