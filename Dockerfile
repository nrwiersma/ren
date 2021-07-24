FROM alpine:latest as builder

FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY ren /ren

ENV PORT "80"

EXPOSE 80
CMD ["./ren", "server"]
