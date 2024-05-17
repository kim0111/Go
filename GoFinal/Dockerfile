## use official Golang image
#FROM golang:1.22.0
#
## set working directory
#WORKDIR /gomidterm
#
## Copy the source code
#COPY . .
#
## Download and install the dependencies
#RUN go get -d -v ./...
#
## Build the Go app
#RUN go build -o applestore ./cmd/apple
#
##EXPOSE the port
#EXPOSE 8081
#
## Run the executable
#CMD ["./applestore"]


# Start from the official Golang image.
FROM golang:1.22.0 as builder

# Set the Current Working Directory inside the container
WORKDIR /gomidterm

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy migration files

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application, disable CGO to create a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o apple-store ./cmd/apple

# Use a smaller image to run the app
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN ls -la

# Copy the pre-built binary file from the previous stage
COPY --from=builder /gomidterm/ .
COPY --from=builder /gomidterm/pkg/apple/migrations ./migrations

# Command to run the executable
CMD ["./apple-store"]