FROM golang as go
WORKDIR /workspace
COPY go* /workspace/
COPY collector /workspace/collector
COPY renderman /workspace/renderman
RUN go build -o /workspace/bin/collector ./collector
RUN go build -o /workspace/bin/renderman ./renderman

FROM browserless/chrome

USER root
RUN apt-get update -y && apt-get install -y --no-install-recommends supervisor tzdata nginx gettext

ENV DISABLED_FEATURES='["pdfEndpoint", "contentEndpoint"]'
ENV DEFAULT_USER_DATA_DIR=/tmp/browserless
ENV ENABLE_CORS=true
ENV ENABLE_DEBUGGER=false
ENV EXIT_ON_HEALTH_FAILURE=true
ENV KEEP_ALIVE=true
ENV MAX_CONCURRENT_SESSIONS=2
ENV MAX_QUEUE_LENGTH=2
ENV PREBOOT_CHROME=true
ENV DEFAULT_IGNORE_HTTPS_ERRORS=true
ENV DEFAULT_LAUNCH_ARGS='["--disable-canvas-aa", "--disable-2d-canvas-clip-aa", "--disable-gl-drawing-for-tests","--disable-setuid-sandbox", "--disable-dev-shm-usage", "--disable-accelerated-2d-canvas","--disable-gpu", "--disable-infobars", "--disable-breakpad", "--use-gl=swiftshader","--hide-scrollbars", "--mute-audio", "--no-first-run", "--no-sandbox", "--window-size=1366,768","--disable-web-security"]'

COPY --from=go   /workspace/bin/collector /workspace/collector
COPY --from=go   /workspace/bin/renderman /workspace/renderman

COPY supervisord.conf /etc/supervisor/supervisord.conf
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf.template

RUN mkdir /var/run/supervisor && mkdir /var/run/nginx \
    && mkdir -p /var/log/nginx \
    && adduser --disabled-password --disabled-login renderman \
    && chown -R renderman:renderman /workspace \
    && chown -R renderman:renderman /var/run/supervisor \
    && chown -R renderman:renderman /var/run/nginx \
    && chown -R renderman:renderman /var/lib/nginx \
    && chown -R renderman:renderman /var/log/nginx \
    && touch /etc/nginx/conf.d/default.conf && chown -R renderman:renderman /etc/nginx/conf.d/default.conf

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log \
    && ln -sf /dev/stderr /var/log/nginx/error.log

USER renderman
EXPOSE 8081 9090 3000
ENTRYPOINT [ "/usr/bin/supervisord" ]
CMD ["-c", "/etc/supervisor/supervisord.conf"]
