.PHONY: all

container-build:
	docker build -t moshi-moshi:${TAG} .

container-start:
	docker-compose up -d

container-stop:
	docker-compose stop
