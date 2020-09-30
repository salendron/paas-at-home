FROM golang:1.14.1-alpine3.11

RUN apk add --no-cache git gcc libc-dev

# Install Libs needed for vscode
RUN go get golang.org/x/tools/gopls
RUN go get github.com/go-delve/delve/cmd/dlv

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download

COPY ./src/. ./src

# Build the Go app
RUN go build -o ./app/server ./src/.

# Run the binary
CMD ["./app/server"]