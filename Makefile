BINARY := timezone-utility

.PHONY: build test vet fmt cross clean

build:
	go build -o bin/$(BINARY) ./cmd/timezone-utility

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

cross:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux-amd64 ./cmd/timezone-utility
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY)-linux-arm64 ./cmd/timezone-utility
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY)-darwin-amd64 ./cmd/timezone-utility
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY)-darwin-arm64 ./cmd/timezone-utility
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY)-windows-amd64.exe ./cmd/timezone-utility
	GOOS=windows GOARCH=arm64 go build -o bin/$(BINARY)-windows-arm64.exe ./cmd/timezone-utility

clean:
	rm -rf bin/
