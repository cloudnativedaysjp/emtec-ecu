CNDCTL_IMAGE_NAME ?= $(REGISTRY_BASE)/emtectl
CNDCTL_IMAGE_TAG ?= develop

##@ Build

.PHONY: emtectl-all
emtectl-all: emtectl-build-image emtectl-push-image ## build Docker image

.PHONY: emtectl-build-image
emtectl-build-image: cmd/emtectl/Dockerfile ## build Docker image
	docker build . -f cmd/emtectl/Dockerfile -t $(CNDCTL_IMAGE_NAME):$(CNDCTL_IMAGE_TAG)

.PHONY: emtectl-push-image
emtectl-push-image: ## push Docker image
	docker push $(CNDCTL_IMAGE_NAME):$(CNDCTL_IMAGE_TAG)

