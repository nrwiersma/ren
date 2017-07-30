# Build container
FROM golang:1.8 as builder

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/nrwiersma/ren/
COPY ./ .
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s' -o ren ./cmd/ren

# Run container
FROM scratch

COPY --from=builder /go/src/github.com/nrwiersma/ren/ren .
COPY ./templates/ ./templates/

ENV REN_PORT "80"

EXPOSE 80
CMD ["./ren", "server"]