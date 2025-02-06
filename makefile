build:
	cd cmd/na-meste-api && go build -o ../../bin/ && cd ../..

run:
	./bin/na-meste-api

migrate:
	./bin/na-meste-api --migrate

clear:
	rm bin/na-meste-api