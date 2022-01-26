FROM golang as builder


WORKDIR /go/src/github.com/user/app
COPY . .
RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


# Docker run Golang app
FROM alpine
ARG LISTEN_PORT=1234
WORKDIR /root/
COPY --from=builder /go/src/github.com/user/app .
EXPOSE $LISTEN_PORT
CMD ["./app"]