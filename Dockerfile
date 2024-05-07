FROM scratch

COPY app /
COPY config.json /config.json

ENTRYPOINT ["/app"]
