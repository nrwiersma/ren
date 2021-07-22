FROM scratch

COPY ren /ren

ENV PORT "80"
ENV TEMPLATES "file:///templates"

EXPOSE 80
CMD ["./ren", "server"]
