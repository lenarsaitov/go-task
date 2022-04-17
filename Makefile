build:
	docker build -t go_server . --target=server

run:
	docker run --rm -p 10000:10000 -p 5440:5432 go_server

migrate:
	goose postgres "host=localhost port=5440 password=docker user=docker dbname=docker sslmode=disable"