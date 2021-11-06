build:
	go build -o storage-challenges ../
	go build -o cmd/start-test ./cmd

migrate:
	STORAGE_CHALLENGE_CONFIG=node0 ./storage-challenges --migrate-seed

start-node0:
	# node0
	STORAGE_CHALLENGE_CONFIG=node0 ./storage-challenges --test
  
start-node1:
	# node1
	STORAGE_CHALLENGE_CONFIG=node1 ./storage-challenges --test

start-node2:
	# node2
	STORAGE_CHALLENGE_CONFIG=node2 ./storage-challenges --test

start-node3:
	# node3
	STORAGE_CHALLENGE_CONFIG=node3 ./storage-challenges --test

start-test:
	# start test
	STORAGE_CHALLENGE_CONFIG=cmd ./cmd/start-test
