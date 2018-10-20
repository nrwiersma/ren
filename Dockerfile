# Build container
FROM golang:1.11 as builder

ENV GO111MODULE=on

WORKDIR /app/
COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -X main.version=$(git describe --tags --always)" -o ren ./cmd/ren

# Run container
FROM scratch

COPY --from=builder /app/ren .

ENV PORT "80"
ENV TEMPLATES "templates"

EXPOSE 80
CMD ["./ren", "server"]
