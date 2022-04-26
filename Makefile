run:
	go run main.go serve

build:
	go build -o bin/newssite

test:
	go test ./...

migrate:
	./migrations/migrate.sh

totallines:
	find . -name '*.go' | xargs wc -l