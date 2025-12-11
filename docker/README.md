# Hive-Pulse Deployment Guide (Docker Compose)

This guide explains how to deploy Hive-Pulse services and dependencies locally using Docker Compose.

> [!Important]
> Don't forget to change ports from Kubernetes-compatible to Docker-compatible.

## 1. Navigate to Docker Folder

```bash
cd docker
```

## 2. Initialize Helper Services

You can start these in any order:

```bash
cd postgresql && docker compose up -d && cd ..
cd kafka && docker compose up -d && cd ..
cd emqx && docker compose up -d && cd ..
cd clickhouse && docker compose up -d && cd ..
cd metrics && docker compose up -d && cd ..
```

> [!Important]
> Each service may take a few seconds to start.

## 3. Initialize Go Services

Start these services **in order**:

```bash
cd auth && docker compose up -d && cd ..
cd ingress && docker compose up -d && cd ..
cd consumer && docker compose up -d && cd ..
```

> [!Important]
> Services depend on helper services, so make sure the previous steps are completed successfully.

## 4. Run tests (in order to registrate devices)

```bash
cd auth/test
python main.py
```

## 5. Access Metrics and Dashboards

1. Open Grafana: http://127.0.0.1:30000/connections/datasources/grafana-clickhouse-datasource
2. Click `Install`
3. Click `Add new data source`
4. Rename the data source to `clickhouse`
5. Set server address: `clickhouse.hive-pulse.svc.cluster.local`
6. Set server port: `9000`
7. Click `Save & test`
8. Import dashboards: http://127.0.0.1:30000/dashboard/import
9. Use [clickhouse.json](https://github.com/DangeL187/HivePulse/blob/main/k8s/metrics/clickhouse.json) as the dashboard file

## 6. Run devices

```bash
cd device
go run cmd/main.go
```
