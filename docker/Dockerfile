FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY brzaguza .

ENV BRZAGUZA_DATA_DIR="/data"
ENV BRZAGUZA_LOG_DIR="/data/log"

ENTRYPOINT [ "./brzaguza" ]

VOLUME [ "/data" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/tminaorg/brzaguza"