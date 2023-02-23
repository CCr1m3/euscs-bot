upgrade-environment:
	docker-compose down -v
	git pull
	docker build -t euscs/euscs-bot .
	docker-compose up -d