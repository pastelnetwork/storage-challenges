FROM golang:1.17-alpine as builder

RUN apk add gcc g++ libc-dev 

WORKDIR /src/storage-challenge

COPY . /src/storage-challenge

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -extldflags "-static"' -a -installsuffix cgo -o storage-challenges .

FROM scratch

# Copy our static executable
COPY --from=builder /src/storage-challenge/storage-challenges /storage-challenges

# Run the hello binary
ENTRYPOINT ["/storage-challenges", "--test"]
