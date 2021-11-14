FROM golang:1.17-alpine as cmd-builder

RUN apk add gcc g++ libc-dev 

WORKDIR /src

COPY . /src

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -extldflags "-static"' -a -installsuffix cgo -o start-test ./test_nodes/cmd

FROM scratch

# Copy our static executable
COPY --from=cmd-builder /src/start-test /src/start-test

# Run the hello binary
ENTRYPOINT ["/src/start-test"]
