SHELL = /bin/bash

REPO = antonosmond
NAME = host-mutator
NAMESPACE = default
TAG ?= latest

.PHONY: test
test:
	go test -v ./...

.PHONY: image
image: test
	docker build \
		-t ${REPO}/${NAME}:${TAG} \
		-t ${REPO}/${NAME}:latest \
		.

.PHONY: release
release: image
	docker push ${REPO}/${NAME}:${TAG}
	docker push ${REPO}/${NAME}:latest

.PHONY: apply
apply:
	kubectl apply -f manifests \
		--namespace ${NAMESPACE}

.PHONY: example
example:
	kubectl label namespace default host-mutator=enabled --overwrite=true
	kubectl apply -f example \
	  --namespace ${NAMESPACE}
