CNDCTL_IMAGE_NAME ?= $(REGISTRY_BASE)/cndctl
CNDCTL_IMAGE_TAG ?= develop

##@ Build

.PHONY: cndctl-all
cndctl-all: cndctl-build-image cndctl-push-image ## build Docker image

.PHONY: cndctl-build-image
cndctl-build-image: cmd/cndctl/Dockerfile ## build Docker image
	docker build . -f cmd/cndctl/Dockerfile -t $(CNDCTL_IMAGE_NAME):$(CNDCTL_IMAGE_TAG)

.PHONY: cndctl-push-image
cndctl-push-image: ## push Docker image
	docker push $(CNDCTL_IMAGE_NAME):$(CNDCTL_IMAGE_TAG)

