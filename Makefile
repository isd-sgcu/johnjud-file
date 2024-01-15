proto:
	go get github.com/isd-sgcu/johnjud-go-proto@latest

publish:
	cat ./token.txt | docker login --username isd-team-sgcu --password-stdin ghcr.io
	docker build . -t ghcr.io/isd-sgcu/johnjud-file
	docker push ghcr.io/isd-sgcu/johnjud-file

mock-gen:
	mockgen -source ./pkg/client/bucket/bucket.client.go -destination ./mocks/client/bucket/bucket.mock.go

test:
	go vet ./...
	go test  -v -coverpkg ./internal/... -coverprofile coverage.out -covermode count ./internal/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	. ./export-env.sh ; go run ./cmd/.
