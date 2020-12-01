FROM golang:alpine as go
WORKDIR /workspace
COPY go* /workspace/
COPY remote    /workspace/remote
COPY collector /workspace/collector
COPY renderman /workspace/renderman
RUN go build -o /workspace/bin/chrome-launcher ./remote
RUN go build -o /workspace/bin/collector ./collector
RUN go build -o /workspace/bin/renderman ./renderman

FROM alpine
RUN apk add --no-cache supervisor chromium tzdata nginx gettext libintl dumb-init

WORKDIR /workspace

COPY --from=go   /workspace/bin/chrome-launcher /workspace/chrome-launcher
COPY --from=go   /workspace/bin/collector /workspace/collector
COPY --from=go   /workspace/bin/renderman /workspace/renderman

# COPY --from=dart /workspace/app /workspace/app
COPY supervisord.conf /etc/supervisor/supervisord.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf.template
RUN echo "pid        /var/run/nginx/nginx.pid;" >> /etc/nginx/nginx.conf
RUN cat /etc/nginx/nginx.conf

RUN mkdir /var/run/supervisor && mkdir /var/run/nginx \
    && mkdir -p /var/log/nginx \
    && adduser -D renderman \
    && chown -R renderman:renderman /workspace \
    && chown -R renderman:renderman /var/run/supervisor \
    && chown -R renderman:renderman /var/run/nginx \
    && chown -R renderman:renderman /var/lib/nginx \
    && chown -R renderman:renderman /var/log/nginx \
    && chown -R renderman:renderman /etc/nginx/conf.d/default.conf

USER renderman
EXPOSE 8081 9090 9222
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
