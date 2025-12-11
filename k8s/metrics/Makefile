RELEASE=metrics
NAMESPACE=hive-pulse

.PHONY: up down start stop status

up:
	@echo "Applying Prometheus manifest..."
	kubectl apply -f ./prometheus.yaml -n $(NAMESPACE)
	@echo "Applying Grafana manifest..."
	kubectl apply -f ./grafana.yaml -n $(NAMESPACE)

down:
	@echo "Deleting Prometheus manifest..."
	kubectl delete -f ./prometheus.yaml -n $(NAMESPACE)
	@echo "Deleting Grafana manifest..."
	kubectl delete -f ./grafana.yaml -n $(NAMESPACE)

start:
	@echo "Starting metrics..."
	kubectl scale deployment/$(RELEASE) --replicas=1 -n $(NAMESPACE)

stop:
	@echo "Stopping metrics..."
	kubectl scale deployment/$(RELEASE) --replicas=0 -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
