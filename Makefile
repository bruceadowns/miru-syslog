.PHONY: all run docker

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
	@-docker rmi docker.phx1.jivehosted.com/miru/miru-syslog 2>/dev/null

build:
	@echo build service
	@go build github.com/bruceadowns/miru-syslog

install:
	@echo install service
	@go install github.com/bruceadowns/miru-syslog

docker:
	@echo build docker image
	@docker build --no-cache --file Dockerfile --tag docker.phx1.jivehosted.com/miru/miru-syslog:latest .

run:
	@echo run docker image
	@docker run -it --rm -p 514:514 -p 514:514/udp -p 8081:8081 --env-file mako_env.sh docker.phx1.jivehosted.com/miru/miru-syslog:latest
