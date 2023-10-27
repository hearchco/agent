FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY brzaguza .

ENV \
    BRZAGUZA_DATA="/config" \
    BRZAGUZA_LOG="/config/log"

ENTRYPOINT [ "./brzaguza" ]

VOLUME [ "/config" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/tminaorg/brzaguza"