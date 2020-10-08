FROM golang:alpine as builder
ENV UID 2000
RUN adduser --uid $UID --disabled-password user
RUN apk add --no-cache git curl bash
COPY main.go .
RUN mkdir /out
EXPOSE 8080
ENV PORT=8080
RUN go build -o /out/proxy
CMD /out/proxy
USER $UID
