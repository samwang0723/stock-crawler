# docker build -t samwang0723/redis-sentinel -f Dockerfile .
FROM amd64/redis:5.0-alpine
MAINTAINER Sam Wang <sam.wang.0723@gmail.com>
EXPOSE 26379
COPY ./sentinel.conf /etc/redis/sentinel.conf
ENTRYPOINT redis-server /etc/redis/sentinel.conf --sentinel
