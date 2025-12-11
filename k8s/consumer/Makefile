IMAGE_NAME=consumer
IMAGE_TAG=latest
NAMESPACE=hive-pulse
REGISTRY=localhost:5000

.PHONY: build up down status

build:
	@echo "Building consumer..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f ../../consumer/Dockerfile ../../consumer/
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

up:
	@echo "Starting consumer..."
	kubectl apply -f ./consumer.yaml -n $(NAMESPACE)

down:
	@echo "Stopping consumer..."
	kubectl delete deployment $(IMAGE_NAME) -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
