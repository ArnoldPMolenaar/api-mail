FROM golang:1.23-alpine
LABEL authors="Arnold Molenaar <arnold.molenaar@webmi.nl> (https://arnoldmolenaar.nl/)"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify
RUN go install github.com/air-verse/air@v1.61.1

# Copy everything from the current directory to the Working Directory inside the container
COPY ./ /app

EXPOSE 5000

# Run the app
CMD ["air"]
