# Copy from https://t.me/golangID/138064
FROM golang:1.18.0-alpine3.15 AS builder

RUN apk update
RUN apk add --no-cache git

WORKDIR /app

# COPY .git ./.git

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /main -ldflags "-X main.buildDate=$(date +%Y%m%d%H%M%S) -X main.commitHash=$(git rev-parse HEAD)"

FROM alpine:3.15.4

WORKDIR /app

COPY --from=builder /main ./main
COPY ./docs ./docs
COPY ./tmpl ./tmpl

ARG PORT=3000
EXPOSE ${PORT}

ENTRYPOINT ["/app/main"]