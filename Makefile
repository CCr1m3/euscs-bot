upgrade-environment:
	git pull
	docker build -t euscs/euscs-bot .
	docker-compose down -v
	docker-compose up -d