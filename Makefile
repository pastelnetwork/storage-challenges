build:
	go generate
	go mod tidy
	go mod vendor
	go build -o ./test_nodes/storage-challenges .
	docker rmi -f sc sc-testnode
	docker build -t sc .
	docker build -t sc-testnode -f testnode.Dockerfile .
	unzip -d test_nodes/sample_raptorq_symbol_files 'test_nodes/sample_raptorq_symbol_files_zip/*.zip'
	unzip -d test_nodes/incremental_raptorq_symbol_files 'test_nodes/incremental_raptorq_symbol_files_zip/*.zip'
	ln -s -f test_nodes/sample_raptorq_symbol_files sample_raptorq_symbol_files
	cp -ru test_nodes/sample_raptorq_symbol_files/* test_nodes/sample_raptorq_symbol_files_29
	cp -ru test_nodes/sample_raptorq_symbol_files/* test_nodes/sample_raptorq_symbol_files_41
	cp -ru test_nodes/sample_raptorq_symbol_files/* test_nodes/sample_raptorq_symbol_files_47
	cp -ru test_nodes/sample_raptorq_symbol_files/* test_nodes/sample_raptorq_symbol_files_87
	cp -ru test_nodes/incremental_raptorq_symbol_files/* test_nodes/incremental_raptorq_symbol_files_29
	cp -ru test_nodes/incremental_raptorq_symbol_files/* test_nodes/incremental_raptorq_symbol_files_41
	cp -ru test_nodes/incremental_raptorq_symbol_files/* test_nodes/incremental_raptorq_symbol_files_47
	cp -ru test_nodes/incremental_raptorq_symbol_files/* test_nodes/incremental_raptorq_symbol_files_87
migrate:
	rm -f ./test_nodes/storage-challenge.sqlite
	STORAGE_CHALLENGE_CONFIG=config ./test_nodes/storage-challenges --migrate-seed
	mv storage-challenge.sqlite test_nodes/

start-nodes:
	docker-compose -f docker-compose.testnodes.yml rm -s -f
	docker-compose -f docker-compose.testnodes.yml up

start-test:
	docker rm -f testnode
	# start test
	docker run --name testnode --env STORAGE_CHALLENGE_CONFIG=/src \
	--network storage_challenges_default \
	--network-alias testnode \
	-v ${PWD}/test_nodes/cmd/config.yaml:/src/config.yaml \
	-v ${PWD}/test_nodes/storage-challenge.sqlite:/storage-challenge.sqlite \
	-v ${PWD}/test_nodes/sample_raptorq_symbol_files:/sample_raptorq_symbol_files \
	-v ${PWD}/test_nodes/incremental_raptorq_symbol_files:/incremental_raptorq_symbol_files \
	sc-testnode
