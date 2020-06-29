FROM golang:alpine as builder
RUN apk add --no-cache git
RUN apk add --no-cache curl
COPY main.go .
RUN mkdir /out
EXPOSE 8080
ENV PORT=8080
RUN go build -o /out/proxy
CMD /out/proxy
