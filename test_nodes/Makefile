build:
	go build -o storage-challenges ../
	go build -o cmd/start-test ./cmd
	docker build -t sc .

migrate:
	STORAGE_CHALLENGE_CONFIG=node0 ./storage-challenges --migrate-seed

start-nodes:
	docker-compose up

start-test:
	# start test
	STORAGE_CHALLENGE_CONFIG=cmd .cmd/start-test
