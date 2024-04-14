build:
	go build -o banner ./cmd/banner

run:
	go run ./...

clean:
	rm banner

up:
	docker compose up --build
