FROM golang:1.23-alpine
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENTRYPOINT ["goose"]