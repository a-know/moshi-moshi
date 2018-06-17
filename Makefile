.PHONY: all

container-build:
	docker build -t moshi-moshi:${TAG} .

build-for-gcr:
	docker build -t asia.gcr.io/${PROJECT}/moshi-moshi:latest .

push-to-gcr:
	gcloud docker -- push asia.gcr.io/${PROJECT}/moshi-moshi:latest

container-start:
	docker-compose up -d

k8s-deploy:
	kubectl create -f deployment.yml

k8s-expose:
	kubectl expose deployment moshi-moshi --type=LoadBalancer --port 80 --target-port=8080

container-stop:
	docker-compose stop
