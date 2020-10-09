FROM golang:alpine as builder
RUN adduser --uid 2000 --disabled-password user
RUN apk add --no-cache git curl bash
COPY main.go .
RUN mkdir /out
EXPOSE 8080
ENV PORT=8080
RUN go build -o /out/proxy
CMD /out/proxy
USER 2000
