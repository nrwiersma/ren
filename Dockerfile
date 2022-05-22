FROM  gcr.io/distroless/static:nonroot

COPY ren /ren

ENV ADDR ":8080"

EXPOSE 8080
CMD ["/ren", "server"]
