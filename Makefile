BINARYFILE=./build/bin/app
SOURCEFILE=./cmd/app/main.go
DOCKER_NAME=computer-club-system

build: $(SOURCEFILE)
	go build -o $(BINARYFILE) $(SOURCEFILE)

build-linux: $(SOURCEFILE)
	GOOS=linux go build -o $(BINARYFILE) $(SOURCEFILE)

run-test_file1: build
	$(BINARYFILE) ./tests/test_file1.txt

docker-build:
	docker build --tag $(DOCKER_NAME) .

docker-run-test_file1: docker-build
	@docker rm -f -v $(DOCKER_NAME)
	docker run --name $(DOCKER_NAME) $(DOCKER_NAME) ./tests/test_file1.txt

clean:
	rm $(BINARYFILE)
