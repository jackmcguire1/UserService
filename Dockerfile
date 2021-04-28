FROM golang:1.16.2-stretch

WORKDIR '/app/'

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/main ./cmd/


# This container exposes port 8080 to the outside world
EXPOSE 7755

# Run the binary program produced by `go install`
CMD ["./out/main"]
