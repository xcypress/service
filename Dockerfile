FROM registry.cn-hangzhou.aliyuncs.com/codoon/docker-golang:latest
MAINTAINER xinhp <xinhp.git@gmail.com>
COPY . /go/src/service-dev
RUN go install service-dev
ENTRYPOINT /go/bin/service-dev
