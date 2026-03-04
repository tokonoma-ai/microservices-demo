PWD = $(shell pwd)

.PHONY: render-kind
render-kind:
	kubectl kustomize overlays/kind

.PHONY: render-eks
render-eks:
	kubectl kustomize overlays/eks
