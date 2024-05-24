run:
	@docker-compose --env-file .env up
start:
	@docker-compose --env-file .env up -d
stop:
	@docker-compose down
log:
	@docker-compose logs -f
build-native-amd:
	@cd src && GOOS=linux GOARCH=amd64 go build -v server.go