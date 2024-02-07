# Build image
# https://hub.docker.com/_/golang
FROM golang:1.22.0-alpine AS builder

WORKDIR /usr/src/app

# Pre-copy/cache go.mod for pre-downloading dependencies and only redownloading
# them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o app .

# Runtime image
FROM alpine:3.19
ARG GROUP_ID=1000
ARG USER_ID=1000
ARG USER_HOME=/home/app

COPY --from=builder /usr/src/app/app /usr/local/bin/app

RUN addgroup -g $GROUP_ID -S app \
  && adduser -u $USER_ID -S app -G app -D -h $USER_HOME

USER app
WORKDIR $USER_HOME

ENTRYPOINT ["app"]
CMD []
