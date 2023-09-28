FROM golang:1.21.1-bookworm

ARG BIN_FOLDER=${BIN_FOLDER}
WORKDIR '/app/'

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/main ./cmd/${BIN_FOLDER}/


# This container exposes port 8080 to the outside world
EXPOSE 7755

# Run the binary program produced by `go install`
CMD ["./out/main"]
