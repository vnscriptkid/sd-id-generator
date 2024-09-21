up:
	docker compose up -d

down:
	docker compose down --volumes --rmi all --remove-orphans

cli:
	docker compose exec mysql mysql -u root --password=root_password -D my_database