.PHONY: all

run:
	go build -o moshi-moshi.exe && ./moshi-moshi.exe

container-build:
	docker build -t moshi-moshi:${TAG} .

container-start:
	docker-compose up -d

container-stop:
	docker-compose stop

build-for-gcr:
	docker build -t asia.gcr.io/${PROJECT}/moshi-moshi:${VERSION} .

push-to-gcr:
	gcloud docker -- push asia.gcr.io/${PROJECT}/moshi-moshi:${VERSION}

k8s-deploy:
	kubectl create -f deployment.yml

k8s-expose:
	kubectl expose deployment moshi-moshi --type=LoadBalancer --port 80 --target-port=8080

update-container:
	kubectl set image deployment/moshi-moshi moshi-moshi=asia.gcr.io/moshi-moshi-3373/moshi-moshi:${VERSION}

gke-login:
	gcloud container clusters get-credentials preemptible-cluster
