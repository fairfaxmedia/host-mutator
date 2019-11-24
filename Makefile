SHELL = /bin/bash

REPO = "antonosmond"
NAME = "host-mutator"
NAMESPACE = "default"
TAG ?= "latest"

.PHONY: test
test:
	go test -v ./...

.PHONY: image
image: test
	docker build -t ${REPO}/${NAME} .
	docker tag ${REPO}/${NAME}:latest ${REPO}/${NAME}:${TAG}
	docker push ${REPO}/${NAME}

.PHONY: ssl
ssl:
	./ssl.sh ${NAME} ${NAMESPACE}

.PHONY: secret
secret: ssl
	@ kubectl delete secret ${NAME}-ssl > /dev/null 2>&1 || true
	kubectl create secret generic ${NAME}-ssl \
		--namespace ${NAMESPACE} \
		--from-file=./ssl/cert.pem \
		--from-file=./ssl/cert.key

.PHONY: apply
apply: image secret
	kubectl apply -f manifests \
		--namespace ${NAMESPACE}
