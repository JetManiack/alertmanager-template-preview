# Deployment to Kubernetes

This document describes how to deploy the Alertmanager Template Preview application to a Kubernetes cluster.

## Using Private GHCR Images

If your GitHub repository is private, the images pushed to **GitHub Container Registry (GHCR)** will also be private. To pull these images into a Kubernetes cluster, you must create an `imagePullSecret`.

### 1. Create a GitHub Personal Access Token (PAT)
- Go to [GitHub Settings > Developer settings > Personal access tokens > Fine-grained tokens](https://github.com/settings/tokens?type=beta).
- Generate a new token with the `read:packages` scope (specifically `Packages: Read-only` for the repositories you need).

### 2. Create a Kubernetes Secret
Use the following command to create a secret in your namespace:

```bash
kubectl create secret docker-registry ghcr-auth \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_GITHUB_TOKEN> \
  --docker-email=<YOUR_EMAIL>
```

### 3. Reference the Secret in your Deployment
In your `deployment.yaml`, add the `imagePullSecrets` field:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alertmanager-template-preview
spec:
  template:
    spec:
      containers:
      - name: app
        image: ghcr.io/<YOUR_GITHUB_ORG>/alertmanager-template-preview:latest
      imagePullSecrets:
      - name: ghcr-auth
```

## Configuration

The application can be configured using environment variables or command-line flags.

### Environment Variables
- `PORT`: The port the server listens on (default: `8080`).
- `PROMETHEUS_URL`: The URL of the Prometheus server to use for queries.

### Example Kubernetes Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: alertmanager-template-preview
spec:
  selector:
    app: alertmanager-template-preview
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```
