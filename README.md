# storage-challenges
Pastel Storage Challenges

## install development code generation tools

```
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install github.com/gogo/protobuf/protoc-gen-gogoslick@latest
    go install github.com/vektra/mockery/v2@latest
```

## before run:

```
    go mod tidy
    go generate
    go build -o storage-challenges .
```

Specify where the config.yaml placed in by set the env STORAGE_CHALLENGE_CONFIG, if not set STORAGE_CHALLENGE_CONFIG, default config.yaml file will be in `./config` folder.

## run node:

```
    ./storage-challenges
```

## start debug nodes

- build:

```
make build
```

- migrate and seeding test data:

```
make migrate
```

- start 12 test nodes:

```
make start-nodes
```
