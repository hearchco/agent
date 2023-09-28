FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY brzaguza .

ENV \
    GIN_MODE="release" \
    BRZAGUZA_CONFIG="/config"

ENTRYPOINT [ "./brzaguza" ]

VOLUME [ "/config" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/tminaorg/brzaguza"