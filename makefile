build:
	cd cmd/na-meste-api && go build -o ../../bin/ && cd ../..

build_migrate:
	cd cmd/na-meste-migrate && go build -o ../../bin/ && cd ../..

run:
	./bin/na-meste-api

migrate:
	./bin/na-meste-migrate

clear:
	rm bin/na-meste-api && rm bin/na-meste-migrate