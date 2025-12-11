RELEASE=emqx
NAMESPACE=hive-pulse

.PHONY: up down start stop status

up:
	@echo "Applying EMQX manifest..."
	kubectl apply -f ./emqx.yaml -n $(NAMESPACE)

down:
	@echo "Deleting EMQX manifest..."
	kubectl delete -f ./emqx.yaml -n $(NAMESPACE)

start:
	@echo "Starting EMQX..."
	kubectl scale deployment/$(RELEASE) --replicas=1 -n $(NAMESPACE)

stop:
	@echo "Stopping EMQX..."
	kubectl scale deployment/$(RELEASE) --replicas=0 -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
