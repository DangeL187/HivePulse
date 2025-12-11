RELEASE=clickhouse
NAMESPACE=hive-pulse

.PHONY: up down start stop status

up:
	@echo "Applying ClickHouse manifest..."
	kubectl -n $(NAMESPACE) delete configmap clickhouse-init-sql
	kubectl -n $(NAMESPACE) create configmap clickhouse-init-sql --from-file=init.sql
	kubectl apply -f ./clickhouse.yaml -n $(NAMESPACE)

down:
	@echo "Deleting ClickHouse manifest..."
	kubectl delete -f ./clickhouse.yaml -n $(NAMESPACE)

start:
	@echo "Starting ClickHouse..."
	kubectl scale deployment/$(RELEASE) --replicas=1 -n $(NAMESPACE)

stop:
	@echo "Stopping ClickHouse..."
	kubectl scale deployment/$(RELEASE) --replicas=0 -n $(NAMESPACE)

status:
	kubectl get pods -n $(NAMESPACE)
	kubectl get svc -n $(NAMESPACE)
