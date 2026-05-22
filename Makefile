run:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.466621192456120

run-eclipse:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.466621192456120 --date=2026-08-12T16:37:00 --track=sun --zoom=65

run-mercure-transit:
	go run cmd/main.go --lat=50.34809655855528 --lon=3.4666211924561203 --date=2019-11-11T12:00:00 --track=sun --zoom=65

ci: vulncheck vet test

vulncheck:
	govulncheck ./...

vet:
	go vet ./...

test:
	go test -race ./...
