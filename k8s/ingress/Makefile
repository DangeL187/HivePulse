IMAGE_NAME=ingress
IMAGE_TAG=latest
NAMESPACE=hive-pulse
REGISTRY=localhost:5000

.PHONY: build up down status

build:
	@echo "Building ingress..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f ../../ingress/Dockerfile ../../ingress/
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

up:
	@echo "Starting ingress..."
	kubectl apply -f ./ingress.yaml -n $(NAMESPACE)

down:
	@echo "Stopping ingress..."
	kubectl delete deployment $(IMAGE_NAME) -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
