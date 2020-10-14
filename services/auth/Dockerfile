FROM golang:1.14.1-alpine3.11

RUN apk add --no-cache git gcc libc-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./src/. ./src

# Build the Go app
RUN go build -o ./app/server ./src/.

# Run the binary
CMD ["./app/server"]