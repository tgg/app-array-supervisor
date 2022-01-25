FROM golang

# Set the Current Working Directory inside the container
ARG LISTEN_PORT
WORKDIR /app
# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

COPY go.mod ./
COPY go.sum ./
RUN go mod download


# Install the package

RUN go build -o ./app-array-supervisor

# This container exposes port 8080 to the outside world
EXPOSE $LISTEN_PORT

# Run the executable
CMD ["./app-array-supervisor"]