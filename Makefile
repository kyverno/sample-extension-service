REGISTRY := ghcr.io/kyverno/sample-extension-service

build:
	go build

ko:
	KO_DOCKER_REPO=$(REGISTRY) ko build --bare