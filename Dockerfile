FROM golang:alpine as go
WORKDIR /home/renderman
ENV GOPATH=/home/renderman/go

COPY go.* /home/renderman/
COPY collector /home/renderman/collector
COPY remote /home/renderman/remote
COPY renderman /home/renderman/renderman

RUN mkdir -p ${GOPATH} \
    && go build -o /home/renderman/bin/collector ./collector \
    && go build -o /home/renderman/bin/remote    ./remote \
    && go build -o /home/renderman/bin/renderman ./renderman

FROM alpine:3.13

ENV HOME /home/renderman
ENV COLLECTOR_SCHEDULER "*/10 * * * *"
WORKDIR /home/renderman

COPY --from=go   /home/renderman/bin/collector      /home/renderman/collector
COPY --from=go   /home/renderman/bin/remote         /home/renderman/remote
COPY --from=go   /home/renderman/bin/renderman      /home/renderman/renderman

RUN adduser -D renderman \
  && echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" > /etc/apk/repositories \
  && echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
  && echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
  && echo "http://dl-cdn.alpinelinux.org/alpine/v3.12/main" >> /etc/apk/repositories \
  && apk upgrade -U -a \
  && apk add --no-cache ca-certificates supervisor chromium tzdata nginx gettext libintl dumb-init \
  && ln -sf /dev/stdout /var/log/nginx/access.log \
  && ln -sf /dev/stderr /var/log/nginx/error.log \
  && rm -rf /var/cache/* \
  && mkdir /var/cache/apk

COPY crontab.txt      /home/renderman/crontab.txt
COPY supervisord.conf /home/renderman/supervisord.conf
COPY local.conf       /etc/fonts/local.conf
COPY nginx.conf       /etc/nginx/nginx.conf
COPY default.conf     /etc/nginx/conf.d/default.conf.template

RUN  mkdir -p /var/run/supervisor \
  && mkdir -p /var/run/nginx \
  && chown -R renderman:renderman /var/run/supervisor \
  && chown -R renderman:renderman /var/run/nginx \
  && chown -R renderman:renderman /var/lib/nginx \
  && chown -R renderman:renderman /var/log/nginx \
  && chown -R renderman:renderman /etc/nginx

USER renderman
EXPOSE 8082 9090
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/supervisord", "-c", "/home/renderman/supervisord.conf"]
