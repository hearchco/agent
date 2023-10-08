FROM gcr.io/distroless/static-debian12:debug-nonroot

WORKDIR /app

COPY brzaguza-bin .

ENV \
    BRZAGUZA_CONFIG="/config"

ENTRYPOINT [ "sh" ]

# -c is needed because entrypoint is sh
CMD [ "-c", "./brzaguza-bin", "-vv" ]

VOLUME [ "/config" ]

EXPOSE 3030

LABEL org.opencontainers.image.source="https://github.com/tminaorg/brzaguza"