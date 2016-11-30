.PHONY: all run docker

GITHUB_TAG = github.com/bruceadowns/miru-syslog
DOCKER_TAG = docker.phx1.jivehosted.com/r2e2/miru-syslog:latest

all: check test install

check: goimports govet

goimports:
	@echo checking go imports...
	@goimports -d .

govet:
	@echo checking go vet...
	@go tool vet .

test:
	@go get
	@go test -v ./...

clean:
	@-rm -v miru-syslog 2>/dev/null
	@-rm -v $(GOPATH)/bin/miru-syslog 2>/dev/null
	@-docker rmi $(DOCKER_TAG) 2>/dev/null

build:
	@echo build service
	@go build $(GITHUB_TAG)

install:
	@echo install service
	@go install $(GITHUB_TAG)

docker:
	@echo build docker image
	@docker build --no-cache --file docker/Dockerfile --tag $(DOCKER_TAG) .

run:
	@echo run docker image
	@docker run -it --rm -p 8514:8514 -p 8514:8514/udp --env-file docker/env.sh $(DOCKER_TAG)

push:
	@echo push docker image
	@docker push $(DOCKER_TAG)
