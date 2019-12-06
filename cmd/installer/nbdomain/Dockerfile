FROM alpine:latest
RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
    git \
    tzdata

WORKDIR /data
COPY . .

VOLUME ["/data/data"]
EXPOSE 8080
CMD ["/data/nbdomain-linux-amd64"]