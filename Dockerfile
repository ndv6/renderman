FROM golang:alpine as go
WORKDIR /workspace
COPY go* /workspace/
COPY remote/ /workspace/
RUN go build -o chrome-launcher .

FROM google/dart as dart
WORKDIR /workspace
COPY pubspec.yaml .
RUN dart pub get
COPY . .
RUN dart pub get --offline
RUN dart compile exe proxy/main.dart -o /workspace/app

FROM alpine
RUN apk add --no-cache supervisor chromium tzdata dumb-init

COPY --from=dart /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libdl.so.2 /lib/x86_64-linux-gnu/libdl.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=dart /lib/x86_64-linux-gnu/libm.so.6 /lib/x86_64-linux-gnu/libm.so.6
COPY --from=dart /lib/x86_64-linux-gnu/librt.so.1 /lib/x86_64-linux-gnu/librt.so.1
COPY --from=dart /lib/x86_64-linux-gnu/libpthread.so.0 /lib/x86_64-linux-gnu/libpthread.so.0
COPY --from=dart /etc/hosts /etc/hosts
COPY --from=dart /etc/nsswitch.conf /etc/nsswitch.conf
COPY --from=dart /etc/resolv.conf /etc/resolv.conf
COPY --from=dart /lib/x86_64-linux-gnu/libnss_dns.so.2 /lib/x86_64-linux-gnu/libnss_dns.so.2
COPY --from=dart /lib/x86_64-linux-gnu/libresolv.so.2 /lib/x86_64-linux-gnu/libresolv.so.2
COPY --from=dart /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=dart /etc/ssl/certs /etc/ssl/certs

WORKDIR /workspace

COPY --from=go   /workspace/chrome-launcher /workspace/chrome-launcher
COPY --from=dart /workspace/app /workspace/app
COPY supervisord.conf /etc/supervisor/supervisord.conf

RUN mkdir /var/run/supervisor \
    && adduser -D renderman \
    && chown -R renderman:renderman /workspace \
    && chown -R renderman:renderman /var/run/supervisor

USER renderman
EXPOSE 8081 9222
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
