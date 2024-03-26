# Use a golang image as the base image
FROM golang:1.19-alpine

# Set necessary environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /app

# Copy the api server go module files into container, because the server depends on them
COPY ./api /api

# Copy the server go module files into container and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Now, copy the source code into the container
COPY . .

WORKDIR /app/server

# Build the Go app
RUN go build -o kurtosis-package-indexer 

# Set the entry point for the container
ENTRYPOINT ["./kurtosis-package-indexer"]