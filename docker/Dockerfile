FROM cgr.dev/chainguard/static:latest

ENV HEARCHCO_DATA_DIR="/data"
ENV HEARCHCO_LOG_DIR="/data/log"

COPY hearchco .

CMD [ "./hearchco" ]

VOLUME [ "/data" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/hearchco/hearchco"
