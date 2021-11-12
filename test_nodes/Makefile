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
	echo NOTE: update hosts file to point node0->node14 to 127.0.0.1 before start test, but for longterm in future start-test will also running in docker container
	STORAGE_CHALLENGE_CONFIG=cmd ./cmd/start-test
