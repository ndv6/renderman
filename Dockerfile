FROM golang:alpine as go
WORKDIR /workspace
COPY go* /workspace/
COPY remote    /workspace/remote
COPY collector /workspace/collector
COPY renderman /workspace/renderman
RUN go build -o /workspace/bin/chrome-launcher ./remote
RUN go build -o /workspace/bin/collector ./collector
RUN go build -o /workspace/bin/renderman ./renderman

FROM google/dart as dart
WORKDIR /workspace
COPY pubspec.yaml .
RUN dart pub get
COPY . .
RUN dart pub get --offline
RUN dart compile exe renderman_dart/main.dart -o ./renderman_dart

FROM alpine
RUN apk add --no-cache supervisor chromium tzdata nginx gettext libintl dumb-init

WORKDIR /workspace

COPY --from=go   /workspace/bin/chrome-launcher /workspace/chrome-launcher
COPY --from=go   /workspace/bin/collector /workspace/collector
COPY --from=go   /workspace/bin/renderman /workspace/renderman

COPY --from=dart /workspace/renderman_dart /workspace/renderman_dart
COPY --from=dart /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/libdl.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=dart /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
COPY --from=dart /lib/x86_64-linux-gnu/librt.so.1 /lib/x86_64-linux-gnu/librt.so.1
COPY --from=dart /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/libpthread.so.0
COPY --from=dart /lib/x86_64-linux-gnu/libnss_dns.so.2 /lib/x86_64-linux-gnu/libnss_dns.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libresolv.so.2 /lib/x86_64-linux-gnu/libresolv.so.2
COPY --from=dart /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=dart /etc/ssl/certs /etc/ssl/certs
# COPY --from=dart /etc/hosts /etc/hosts
# COPY --from=dart /etc/nsswitch.conf /etc/nsswitch.conf
# COPY --from=dart /etc/resolv.conf /etc/resolv.conf

COPY supervisord.conf /etc/supervisor/supervisord.conf
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf.template

RUN mkdir /var/run/supervisor && mkdir /var/run/nginx \
    && mkdir -p /var/log/nginx \
    && adduser -D renderman \
    && chown -R renderman:renderman /workspace \
    && chown -R renderman:renderman /var/run/supervisor \
    && chown -R renderman:renderman /var/run/nginx \
    && chown -R renderman:renderman /var/lib/nginx \
    && chown -R renderman:renderman /var/log/nginx \
    && chown -R renderman:renderman /etc/nginx/conf.d/default.conf

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

USER renderman
EXPOSE 8081 9090 9222
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
