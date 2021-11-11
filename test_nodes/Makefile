build:
	go build -o storage-challenges ../
	go build -o cmd/start-test ./cmd
	unzip -d sample_raptorq_symbol_files 'sample_raptorq_symbol_files_zip/*.zip'
	unzip -d incremental_raptorq_symbol_files 'incremental_raptorq_symbol_files_zip/*.zip'
	docker build -t sc .

migrate:
	STORAGE_CHALLENGE_CONFIG=node0 ./storage-challenges --migrate-seed

start-nodes:
	docker-compose up

start-test:
	# start test
	STORAGE_CHALLENGE_CONFIG=cmd .cmd/start-test
