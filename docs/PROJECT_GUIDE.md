# Texas Hold'em Project - Complete Step-by-Step Guide

This guide walks you from local development through deployment on Google Kubernetes Engine (GKE) with AMD64 nodes.

---

## Phase 1: Local Development Setup

### 1.1 Install Prerequisites

```bash
# Go (1.21+)
# Arch: sudo pacman -S go
# Or: https://go.dev/dl/

# Flutter
# https://docs.flutter.dev/get-started/install

# Docker
# https://docs.docker.com/get-docker/

# kubectl
# https://kubernetes.io/docs/tasks/tools/
```

### 1.2 Run Backend Locally

```bash
cd backend
go mod download
go run ./cmd/server
# API runs at http://localhost:8080
```

### 1.3 Run Frontend Locally

```bash
cd frontend
flutter pub get
# If web support is missing, run: flutter create . --platforms=web
flutter run -d chrome
```

> **Note:** If `flutter run -d chrome` fails (e.g. missing web files), run `flutter create . --platforms=web` inside the frontend directory to generate the web bootstrap files.

---

## Phase 2: Build Docker Images (AMD64)

GKE nodes are typically AMD64. Always build for that platform:

```bash
# Backend
docker build --platform linux/amd64 -t texas-holdem-backend:latest ./backend

# Frontend (requires Flutter build first)
cd frontend
flutter pub get
flutter build web
docker build --platform linux/amd64 -t texas-holdem-frontend:latest .
cd ..
```

---

## Phase 3: Push to Container Registry

### Option A: Google Container Registry (gcr.io)

```bash
# Authenticate
gcloud auth configure-docker

# Tag for GCR
docker tag texas-holdem-backend:latest gcr.io/YOUR_PROJECT_ID/texas-holdem-backend:latest
docker tag texas-holdem-frontend:latest gcr.io/YOUR_PROJECT_ID/texas-holdem-frontend:latest

# Push
docker push gcr.io/YOUR_PROJECT_ID/texas-holdem-backend:latest
docker push gcr.io/YOUR_PROJECT_ID/texas-holdem-frontend:latest
```

### Option B: Artifact Registry (recommended)

```bash
# Create repository
gcloud artifacts repositories create texas-holdem --repository-format=docker \
  --location=YOUR_REGION --description="Texas Holdem containers"

# Configure docker
gcloud auth configure-docker YOUR_REGION-docker.pkg.dev

# Tag and push
docker tag texas-holdem-backend:latest YOUR_REGION-docker.pkg.dev/YOUR_PROJECT_ID/texas-holdem/backend:latest
docker tag texas-holdem-frontend:latest YOUR_REGION-docker.pkg.dev/YOUR_PROJECT_ID/texas-holdem/frontend:latest

docker push YOUR_REGION-docker.pkg.dev/YOUR_PROJECT_ID/texas-holdem/backend:latest
docker push YOUR_REGION-docker.pkg.dev/YOUR_PROJECT_ID/texas-holdem/frontend:latest
```

---

## Phase 4: Create GKE Cluster (AMD64)

```bash
# Set variables
export PROJECT_ID=your-gcp-project-id
export REGION=us-central1  # or europe-west1, etc.
export CLUSTER_NAME=texas-holdem-cluster

# Create cluster with AMD64 (default, but explicit)
gcloud container clusters create $CLUSTER_NAME \
  --project=$PROJECT_ID \
  --region=$REGION \
  --machine-type=e2-medium \
  --num-nodes=2 \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=5

# Get credentials
gcloud container clusters get-credentials $CLUSTER_NAME \
  --region=$REGION \
  --project=$PROJECT_ID
```

---

## Phase 5: Deploy to Kubernetes

```bash
# Update image URLs in k8s/*.yaml to match your registry
# Then apply:
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml  # Optional: for external access

# Check status
kubectl get pods -n texas-holdem
kubectl get services -n texas-holdem
```

---

## Phase 6: Load Testing

```bash
cd load-test
# Install k6: https://k6.io/docs/getting-started/installation/
k6 run load-test.js
```

---

## Summary Checklist

- [ ] Go backend running locally
- [ ] Flutter frontend running locally
- [ ] Docker images built for linux/amd64
- [ ] Images pushed to GCR or Artifact Registry
- [ ] GKE cluster created
- [ ] kubectl configured
- [ ] Kubernetes manifests applied
- [ ] Load tests passing
