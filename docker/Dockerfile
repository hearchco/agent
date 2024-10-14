FROM scratch

COPY --from=alpine /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=alpine /usr/share/zoneinfo/ /usr/share/zoneinfo/

COPY bootstrap ./hearchco

ENV TZ="UTC"

# "nobody" user
USER 65534

ENTRYPOINT ["./hearchco"]

EXPOSE 8000

LABEL org.opencontainers.image.source="https://github.com/hearchco/agent"
