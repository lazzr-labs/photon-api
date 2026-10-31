# docker
docker-login:
	docker login ghcr.io

docker-logout:
	docker logout ghcr.io

docker-network:
	docker network create docker-network --driver=bridge --attachable

docker-network-stop:
	docker network rm docker-network

docker-postgres:
	docker compose -f docker-postgres.yaml up -d

docker-postgres-stop:
	docker compose -f docker-postgres.yaml down

docker-valkey:
	docker compose -f docker-valkey.yaml up -d

docker-valkey-stop:
	docker compose -f docker-valkey.yaml down

docker-api:
	docker compose -f docker-api.yaml up

docker-api-stop:
	docker compose -f docker-api.yaml down

docker-run:
	$(MAKE) docker-network
	$(MAKE) docker-postgres
	$(MAKE) docker-valkey
	$(MAKE) docker-api

docker-stop:
	$(MAKE) docker-api-stop
	$(MAKE) docker-valkey-stop
	$(MAKE) docker-postgres-stop
	$(MAKE) docker-network-stop

docker-tag:
	sh docker_tag.sh

docker-buildx:
	sh docker_buildx.sh

# helm
helm-install-dev:
	helm install api helm-v1 --values helm-v1/values.yaml -f helm-v1/dev.yaml

helm-upgrade-dev:
	helm upgrade api helm-v1 --values helm-v1/values.yaml -f helm-v1/dev.yaml

helm-install-stag:
	helm install api helm-v1 --values helm-v1/values.yaml -f helm-v1/stag.yaml

helm-upgrade-stag:
	helm upgrade api helm-v1 --values helm-v1/values.yaml -f helm-v1/stag.yaml

helm-install-prod:
	helm install api helm-v1 --values helm-v1/values.yaml -f helm-v1/prod.yaml

helm-upgrade-prod:
	helm upgrade api helm-v1 --values helm-v1/values.yaml -f helm-v1/prod.yaml

helm-uninstall:
	helm uninstall api

# gcp
gcp-auth:
	gcloud auth login

gcp-auth-docker:
	gcloud auth configure-docker us-central1-docker.pkg.dev

gcp-auth-terraform:
	gcloud auth application-default login

gcp-create-vpc:
	gcloud compute networks create vpc-v1 --subnet-mode=auto --project [cloud]-gcp

gcp-delete-vpc:
	gcloud compute networks delete vpc-v1 --quiet --project [cloud]-gcp

gcp-auth-gke:
	gcloud container clusters get-credentials [cloud]-gke-v1-standard --zone us-central1-a --project [cloud]-gcp

# argocd
argocd-create:
	kubectl create namespace [project]-api-stag
	kubectl create namespace [project]-api-prod
	kubectl create namespace argocd
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

argocd-port:
	kubectl port-forward svc/argocd-server -n argocd 8080:443

argocd-password:
	argocd admin initial-password -n argocd
