//go:generate protoc --gogoslick_out=./application/dto --gogoslick_opt=paths=source_relative -I ./proto ./proto/storage_challenges_dto.proto
//go:generate protoc --go-grpc_out=./application/grpc --go-grpc_opt=paths=source_relative -I ./proto ./proto/storage_challenges.proto
//go:generate mockery --name=Client --dir=../gonode/pastel --output=test_nodes/mocks
package main
