build:
	@ go build -o bin/dist-store
run: build
	@ ./bin/dist-store
test:
	@ go test ./... -v
