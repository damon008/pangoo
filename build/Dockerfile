#FROM golang:1.20 as builder
#
#ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,direct

FROM ubuntu:16.04 as builder

RUN apt-get update && apt-get install -y --no-install-recommends \
        g++ \
        ca-certificates \
        wget && \
    rm -rf /var/lib/apt/lists/*


ENV GOLANG_VERSION 1.20.4
RUN wget -nv -O - https://studygolang.com/dl/golang/go1.20.4.linux-amd64.tar.gz \
     | tar -C /usr/local -xz

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV GOPATH /go/pkg
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN apt-get update && apt-get install -y git

COPY . /go/src/pangoo
WORKDIR /go/src/pangoo

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app-katalyst ./katalyst

#go build -o ./bin/katalyst ./katalyst

FROM debian:stretch-slim
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /data && \
    chown -R nobody:nogroup /data

COPY --from=builder /go/src/pangoo/app-katalyst /data
COPY --from=builder /go/src/pangoo/conf /data/conf

RUN rm -rf /go/src/pangoo

USER        nobody
WORKDIR     /data
ENV         CONFIG_ENV=dev
#ENV         HERTZ_CONF_DIR=conf PSM=cwg.bizdemo.bookinfo HERTZ_LOG_DIR=/data/katalyst/logs/
#CMD ["sleep 10000"]
#CMD ["./app-katalyst"]
ENTRYPOINT ["sh","-c","./app-katalyst"]

#docker run --name katalyst -d -p 20000:20000 -v /data/pangoo/conf:/data/conf  -v  /data/katalyst:/data/katalyst katalyst