proto:
	go get github.com/isd-sgcu/johnjud-go-proto@latest

publish:
	cat ./token.txt | docker login --username isd-team-sgcu --password-stdin ghcr.io
	docker build . -t ghcr.io/isd-sgcu/johnjud-file
	docker push ghcr.io/isd-sgcu/johnjud-file

mock-gen:
	mockgen -source ./pkg/repository/image/image.repository.go -destination ./mocks/repository/image/image.mock.go

test:
	go vet ./...
	go test  -v -coverpkg ./src/app/... -coverprofile coverage.out -covermode count ./src/app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./src/.
