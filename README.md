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

```
cd test_nodes
```

- update pastel client mock in test_nodes/mock_pastel_client.go to matches your test cases
- update database dummy data in test_nodes/migations.go to matches your test cases

- build:

```
make build
```

- migrate and seeding test data:

```
make migrate
```

- start test node0:

```
make start-node0
```
- start test node1 in new terminal

```
make start-node1
```
- start test node2 in new terminal

```
make start-node2
```

- start test node3 in new terminal

```
make start-node3
```

- start test process (call to node0 generate storage challenges)

```
make start-test
```
