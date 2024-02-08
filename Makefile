run:
	@docker-compose --env-file .env up
start:
	@docker-compose --env-file .env up -d
log:
	@docker-compose logs -f