FROM golang:alpine as builder

RUN apk update && apk add --no-cache git curl

WORKDIR /app

## Install go-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz 

COPY go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY ../../ .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin ./cmd/api

# Start fresh from a smaller image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/cfg ./cfg
COPY --from=builder /app/.env . 
COPY --from=builder /app/deploy/api/start.sh .
COPY --from=builder /app/internal/db/migrations ./migrations

COPY --from=builder /app/bin/api .

EXPOSE 5000

CMD ["./api"]
ENTRYPOINT ["./start.sh"]