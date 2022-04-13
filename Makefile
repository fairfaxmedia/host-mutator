SHELL = /bin/bash

REPO = fairfaxmedia
NAME = host-mutator
NAMESPACE = infrastructure-ingress-host-mutator
TAG ?= latest
KUBECTL ?= kubectl # --as pe-admin

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
	${KUBECTL} create namespace ${NAMESPACE} --dry-run=client -o yaml | ${KUBECTL} apply -f -
	${KUBECTL} apply --filename manifests --namespace ${NAMESPACE}

.PHONY: example
example:
	${KUBECTL} label namespace ${NAMESPACE} host-mutator=enabled --overwrite=true
	${KUBECTL} apply --filename example --namespace ${NAMESPACE}
