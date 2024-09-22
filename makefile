up:
	docker compose up -d

down:
	docker compose down --volumes --rmi all --remove-orphans

mysql:
	docker compose exec mysql mysql -u root --password=root_password -D my_database

mongo:
	docker compose exec mongo mongosh --host localhost:27017

psql:
	docker compose exec pg psql -U postgres -d postgres
