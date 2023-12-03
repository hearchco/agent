FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY hearchco .

ENV HEARCHCO_DATA_DIR="/data"
ENV HEARCHCO_LOG_DIR="/data/log"

ENTRYPOINT [ "./hearchco" ]

VOLUME [ "/data" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/hearchco/hearchco"