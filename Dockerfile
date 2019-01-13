FROM golang:alpine AS binarybuilder
# Install build deps
RUN apk --no-cache --no-progress add --virtual build-deps build-base git linux-pam-dev
WORKDIR /go/src/github.com/naiba/domain-panel/
COPY . .
RUN cd cmd/domain-panel \
    && go build -ldflags="-s -w"

FROM alpine:latest
RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
    git \
    tzdata
# Copy binary to container
WORKDIR /data
COPY cmd/domain-panel/theme ./theme
COPY --from=binarybuilder /go/src/github.com/naiba/domain-panel/cmd/domain-panel/domain-panel ./domain-panel

# Configure Docker Container
VOLUME ["/data/data"]
EXPOSE 8080
CMD ["/data/domain-panel"]