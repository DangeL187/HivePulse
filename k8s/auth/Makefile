IMAGE_NAME=auth
IMAGE_TAG=latest
NAMESPACE=hive-pulse
REGISTRY=localhost:5000

.PHONY: build up down status

build:
	@echo "Building auth..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f ../../auth/Dockerfile ../../auth/
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

up:
	@echo "Starting auth..."
	kubectl apply -f ./auth.yaml -n $(NAMESPACE)

down:
	@echo "Stopping auth..."
	kubectl delete deployment $(IMAGE_NAME) -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
