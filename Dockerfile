# Container used to build the Go application. This container is large and 
# consumes considerable resources. Therefore, it is not advisable to run the Go
# application from here.
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server/

# This container is small and consumes inconsiderable resources. Therefore, it
# is advisable to run the Go application from here. Only the necessary files
# to run the application, namely ./assets, .env, and the binary, are copied
# from the builder container.
FROM alpine
WORKDIR /app
COPY ./assets ./assets
COPY .env .env
COPY --from=builder /app/server ./server
CMD ./server