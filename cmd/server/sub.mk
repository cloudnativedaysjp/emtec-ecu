SERVER_IMAGE_NAME ?= $(REGISTRY_BASE)/server
SERVER_IMAGE_TAG ?= develop

##@ Build

.PHONY: server-all
server-all: server-build-image server-push-image ## build Docker image

.PHONY: server-build-image
server-build-image: cmd/server/Dockerfile ## build Docker image
	docker build . -f cmd/server/Dockerfile -t $(SERVER_IMAGE_NAME):$(SERVER_IMAGE_TAG)

.PHONY: server-push-image
server-push-image: ## push Docker image
	docker push $(SERVER_IMAGE_NAME):$(SERVER_IMAGE_TAG)

