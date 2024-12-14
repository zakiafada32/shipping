GO_VERSION := 1.23.4  

.PHONY: install-go init-go

setup: install-go init-go 

install-go: 
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz

init-go: 
    echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
    echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.bashrc

upgrade-go: 
	sudo rm -rf /usr/bin/go
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz

build:
	go build -o api cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./... -coverprofile=coverage.out

coverage:
	go tool cover -func coverage.out

min-coverage:
	go tool cover -func coverage.out | grep "total:" | awk '{print ((int($$3) >= 80) != 1) }'

report:
	go tool cover -html=coverage.out -o cover.html

check-format:
	test -z $$(go fmt ./...)

install-lint:
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.2

static-check:
	golangci-lint run

copy-hooks:
	chmod +x scripts/hooks/* 
	cp -r ./scripts/hooks .git/.

install-redis:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm install redis-cluster bitnami/redis --set password=$$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')