
test:
	go test -coverprofile=coverage.out ./ -v
	go tool cover -html=coverage.out -o coverage.html