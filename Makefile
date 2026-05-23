VERSION := $(shell git tag --sort=-v:refname | head -n 1)
MAJOR := $(shell echo $(VERSION) | cut -d. -f1)
MINOR := $(shell echo $(VERSION) | cut -d. -f1,2)

run:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.466621192456120

run-eclipse:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.466621192456120 --date=2026-08-12T16:37:00 --track=sun --zoom=65

run-mercure-transit:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.4666211924561203 --date=2019-11-11T12:00:00 --track=sun --zoom=65

ci: vulncheck mod-verify vet staticcheck lint test

vulncheck:
	govulncheck ./...

mod-verify:
	go mod verify

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run

test:
	go test -race ./...

build:
	docker build --build-arg VERSION=$(VERSION) -t mickaelblondeau/asterminal:$(VERSION) -t mickaelblondeau/asterminal:$(MINOR) -t mickaelblondeau/asterminal:$(MAJOR) -t mickaelblondeau/asterminal:latest .
	docker push mickaelblondeau/asterminal --all-tags
