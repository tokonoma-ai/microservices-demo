.PHONY: render-kind
render-kind:
	make -C deploy/kubernetes render-kind

.PHONY: render-eks
render-eks:
	make -C deploy/kubernetes render-eks
