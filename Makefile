build-and-run:
	go build -o ./bin/entities ./cmd/entities && ./bin/entities
build-and-run-migrate:
	go build -o ./bin/entities ./cmd/entities && ./bin/entities -migrateDB=true
build-and-run-background:
	go build -o ./bin/entities ./cmd/entities && nohup ./bin/entities > /dev/null 2>&1&
build-and-run-docker:
	docker run -p 8080:8080 --volume $$PWD:/app --rm -it $$(docker build -q .)
unit-test:
	 go test -coverprofile coverage.out $$(go list ./... | grep -v integration)
integration-test:
	BASE_URL=http://localhost:8080 go test ./integration -count=1
format:
	go fmt ./...
vet:
	go vet
tidy:
	go mod tidy
clean:
	go ​clean -modcache
find-process:
	pgrep entities
kill-process:
	kill 12345678
find-go-symlink:
	readlink -f /usr/bin/go
mocks:
	mockery
digital-ocean:
	doctl compute ssh vcp-nyc
coverage:
	@go tool cover -html coverage.out -o coverage.html
	explorer.exe coverage.html
