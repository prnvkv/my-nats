PROJECT_ROOT  := github.com/prnvkv/my-nats
BUILD_PATH := bin
DOCKERFILE_PUB := $(CURDIR)/docker/pub/.
DOCKERFILE_SUB := $(CURDIR)/docker/sub/.
DOCKERFILE_QPUB := $(CURDIR)/docker/qpub/.
DOCKERFILE_QSUB := $(CURDIR)/docker/qsub/.
DOCKERFILE_REQ := $(CURDIR)/docker/req/.
DOCKERFILE_REP := $(CURDIR)/docker/rep/.

# configuration for the build path for each of the app
PUBSUB_PUB := $(BUILD_PATH)/pub-sub/pub
PUBSUB_SUB := $(BUILD_PATH)/pub-sub/sub
QGROUP_PUB := $(BUILD_PATH)/queue-group/producer
QGROUP_SUB := $(BUILD_PATH)/queue-group/consumer
REQREPLY_PUB := $(BUILD_PATH)/request-reply/pub
REQREPLY_SUB := $(BUILD_PATH)/request-reply/sub

# configuration for building the docker images
SRCROOT_ON_HOST := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_ON_CONTAINER := /go/src/$(PROJECT_ROOT)
DOCKER_RUNNER := docker run --rm -v $(SRCROOT_ON_HOST):$(SRCROOT_ON_CONTAINER) -w $(SRCROOT_ON_CONTAINER)

# configuration to build go app on host machine
GO_CACHE := -pkgdir $(BUILD_PATH)/go-cache
GO_BUILD_FLAGS ?= $(GO_CACHE) -i -v 
GO_PACKAGES := $(shell go list ./... | grep -v vendor)
BUILD_IMAGE := golang:1.14.4-alpine3.12
GO_BUILDER := $(DOCKER_RUNNER) -e CGO_ENABLED=0 -e GO111MODULE=off
GO_BUILD := $(GO_BUILDER) $(BUILD_IMAGE) go build $(GO_BUILD_FLAGS)

# image configuration for each of the go app 
IMAGE_REGISTRY := pranavkv
IMAGE_VERSION := latest

# pub image
IMAGE_NAME_PUB := dockerpub
IMAGE_FULL_PUB := $(IMAGE_REGISTRY)/$(IMAGE_NAME_PUB):$(IMAGE_VERSION)

# sub image
IMAGE_NAME_SUB := dockersub
IMAGE_FULL_SUB := $(IMAGE_REGISTRY)/$(IMAGE_NAME_SUB):$(IMAGE_VERSION)

# qpub image
IMAGE_NAME_QPUB := dockerqpub
IMAGE_FULL_QPUB := $(IMAGE_REGISTRY)/$(IMAGE_NAME_QPUB):$(IMAGE_VERSION)

# qsub image
IMAGE_NAME_QSUB := dockerqsub
IMAGE_FULL_QSUB := $(IMAGE_REGISTRY)/$(IMAGE_NAME_QSUB):$(IMAGE_VERSION)

# request image
IMAGE_NAME_REQ := dockerreq
IMAGE_FULL_REQ := $(IMAGE_REGISTRY)/$(IMAGE_NAME_REQ):$(IMAGE_VERSION)

# reply image
IMAGE_NAME_REP := dockerrep
IMAGE_FULL_REP := $(IMAGE_REGISTRY)/$(IMAGE_NAME_REP):$(IMAGE_VERSION)

# go app paths
SERVER_PUB := $(PROJECT_ROOT)/cmd/pub-sub/pub/
SERVER_SUB := $(PROJECT_ROOT)/cmd/pub-sub/sub/
SERVER_QPUB := $(PROJECT_ROOT)/cmd/queue-group/producer/
SERVER_QSUB := $(PROJECT_ROOT)/cmd/queue-group/consumer/
SERVER_REQ := $(PROJECT_ROOT)/cmd/request-reply/pub/
SERVER_REP := $(PROJECT_ROOT)/cmd/request-reply/sub/

# Targets to build go app binary

.PHONY: build-pub
build-pub:
	@mkdir -p $(BUILD_PATH)/pub-sub/pub
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(PUBSUB_PUB) $(SERVER_PUB)

.PHONY: build-sub
build-sub:
	@mkdir -p $(BUILD_PATH)/pub-sub/sub
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(PUBSUB_SUB) $(SERVER_SUB)

.PHONY: build-qpub
build-qpub:
	@mkdir -p $(BUILD_PATH)/queue-group/producer
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(QGROUP_PUB) $(SERVER_QPUB)

.PHONY: build-qsub
build-qsub:
	@mkdir -p $(BUILD_PATH)/queue-group/consumer
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(PUBSUB_QSUB) $(SERVER_QSUB)

.PHONY: build-req
build-req:
	@mkdir -p $(BUILD_PATH)/request-reply/pub
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(REQREPLY_PUB) $(SERVER_REQ)

.PHONY: build-rep
build-rep:
	@mkdir -p $(BUILD_PATH)/request-reply/pub
	@echo $(GO_BUILD)
	@$(GO_BUILD) -o $(REQREPLY_SUB) $(SERVER_REP)


# Targets to run docker build for each go app

.PHONY: docker-pub
docker-pub: 
	@docker build -f $(DOCKERFILE_PUB) -t $(IMAGE_FULL_PUB) .
	@docker image prune -f --filter label=stage=server-intermediate

.PHONY: docker-sub
docker-sub: 
	@docker build -f $(DOCKERFILE_SUB) -t $(IMAGE_FULL_SUB) .
	@docker image prune -f --filter label=stage=server-intermediate

.PHONY: docker-qpub
docker-qpub: 
	@docker build -f $(DOCKERFILE_QPUB) -t $(IMAGE_FULL_QPUB) .

.PHONY: docker-qsub
docker-qsub: 
	@docker build -f $(DOCKERFILE_QSUB) -t $(IMAGE_FULL_QSUB) .

.PHONY: docker-req
docker-req: 
	@docker build -f $(DOCKERFILE_REQ) -t $(IMAGE_FULL_REQ) .

.PHONY: docker-rep
docker-rep: 
	@docker build -f $(DOCKERFILE_REP) -t $(IMAGE_FULL_REP) .


# Targets to push docker image of the go app 

.PHONY: push-pub
push-pub: 
	@docker push $(IMAGE_NAME_PUB)

.PHONY: push-sub
push-sub: 
	@docker push $(IMAGE_NAME_SUB)

.PHONY: push-qpub
push-qpub: 
	@docker push $(IMAGE_NAME_QPUB)

.PHONY: push-qsub
push-qsub: 
	@docker push $(IMAGE_NAME_QSUB)

.PHONY: push-req
push-req: 
	@docker push $(IMAGE_NAME_REQ)

.PHONY: push-rep
push-rep: 
	@docker push $(IMAGE_NAME_REP)
