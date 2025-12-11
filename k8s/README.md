# Hive-Pulse Deployment Guide

This guide explains how to deploy Hive-Pulse services and dependencies locally using Kubernetes and Docker.

> [!Note]
> Use [k9s](https://github.com/derailed/k9s) for easier pod monitoring and management.

## 1. Create Namespace

Create the hive-pulse namespace in your Kubernetes cluster:

```bash
kubectl create namespace hive-pulse
```

## 2. Start Local Docker Registry

Start a local Docker registry to store images:

```bash
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

## 3. Navigate to Kubernetes Manifests

```bash
cd k8s
```

## 4. Initialize Helper Services

You can start these in any order:

```bash
cd postgresql && make up && cd ..
cd kafka && make up && cd ..
cd emqx && make up && cd ..
cd clickhouse && make up && cd ..
cd metrics && make up && cd ..
```

> [!Important]
> Each service may take a few seconds to start.

## 5. Initialize Go Services

Start these services **in order**:

```bash
cd auth && make build && make up && cd ..
cd ingress && make build && make up && cd ..
cd consumer && make build && make up && cd ..
```

> [!Important]
> Services depend on helper services, so make sure the previous steps are completed successfully.

## 6. Run tests (in order to registrate devices)

```bash
cd auth/test
python main.py
```

## 7. Access Metrics and Dashboards

1. Open Grafana: http://127.0.0.1:30000/connections/datasources/grafana-clickhouse-datasource
2. Click `Install`
3. Click `Add new data source`
4. Rename the data source to `clickhouse`
5. Set server address: `clickhouse.hive-pulse.svc.cluster.local`
6. Set server port: `9000`
7. Click `Save & test`
8. Import dashboards: http://127.0.0.1:30000/dashboard/import
9. Use [clickhouse.json](https://github.com/DangeL187/HivePulse/blob/main/k8s/metrics/clickhouse.json) as the dashboard file

## 8. Run devices

```bash
cd device
go run cmd/main.go
```
