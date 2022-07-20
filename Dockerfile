FROM golang:1.18.4 as builder
LABEL maintainer="jacky01.zhang@outlook.com"

RUN mkdir -p /go/src/github.com/jackyzhangfudan/utilitycontainer
WORKDIR /go/src/github.com/jackyzhangfudan/utilitycontainer/
COPY . /go/src/github.com/jackyzhangfudan/utilitycontainer/
RUN GOOS=linux go build -a -o utility .

FROM alpine:latest
LABEL maintainer="jacky01.zhang@outlook.com"

WORKDIR /
COPY --from=builder /go/src/github.com/jackyzhangfudan/utilitycontainer/utility .
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8080

CMD ["./utility"]