DOCKERHUB = 192.168.27.146
PROJECTNAME = mongo-operator
IMAGE = $(DOCKERHUB)/$(PROJECTNAME)
VERSION = latest

.PHONY: all
all: gen-k8s update build push

.PHONY: gen-k8s
gen-k8s:
	operator-sdk generate k8s --verbose

.PHONY: update
update:
	operator-sdk build $(IMAGE):$(VERSION)

.PHONY: build
build: gen-k8s update

.PHONY: push
push:
	docker push $(IMAGE):$(VERSION)

.PHONY: run
run: export OPERATOR_NAME  = mongo-operator

run:
	operator-sdk up local --namespace=default --operator-flags "--zap-devel=true"

.PHONY: test
test:
	@echo PLS TEST IT.

