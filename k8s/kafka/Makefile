RELEASE=kafka
NAMESPACE=hive-pulse

.PHONY: up down start stop status

up:
	@echo "Applying Kafka manifest..."
	kubectl apply -f ./kafka.yaml -n $(NAMESPACE)

down:
	@echo "Deleting Kafka manifest..."
	kubectl delete -f ./kafka.yaml -n $(NAMESPACE)

start:
	@echo "Starting Kafka..."
	kubectl scale deployment/$(RELEASE) --replicas=1 -n $(NAMESPACE)

stop:
	@echo "Stopping Kafka..."
	kubectl scale deployment/$(RELEASE) --replicas=0 -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
