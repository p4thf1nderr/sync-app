VERSION=0.0.1

start:
	go run cmd/syncd/main.go folders/origin folders/copy

test:
	go test ./... -coverprofile=profile.out && go tool cover -func=profile.out

bench:
	go test ./... -bench=. -benchmem
