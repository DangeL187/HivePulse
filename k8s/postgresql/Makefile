RELEASE=postgresql
NAMESPACE=hive-pulse

.PHONY: up down start stop status

up:
	@echo "Applying PostgreSQL manifest..."
	kubectl -n $(NAMESPACE) delete configmap postgres-init-sql
	kubectl -n $(NAMESPACE) create configmap postgres-init-sql --from-file=init.sql
	kubectl apply -f ./postgresql.yaml -n $(NAMESPACE)

down:
	@echo "Deleting PostgreSQL manifest..."
	kubectl delete -f ./postgresql.yaml -n $(NAMESPACE)

start:
	@echo "Starting PostgreSQL..."
	kubectl scale deployment/$(RELEASE) --replicas=1 -n $(NAMESPACE)

stop:
	@echo "Stopping PostgreSQL..."
	kubectl scale deployment/$(RELEASE) --replicas=0 -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
