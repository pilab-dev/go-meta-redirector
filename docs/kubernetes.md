# Deploying GoMeta Redirector on Kubernetes

This guide walks you through deploying GoMeta Redirector on Kubernetes.

## Prerequisites

- Kubernetes cluster (GKE, EKS, AKS, or local with minikube/kind)
- kubectl configured
- Docker image pushed to `ghcr.io/pilab-dev/go-meta-redirector:latest`
- Domain names pointing to your cluster's ingress

## Quick Deploy

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n go-meta-redirector
```

## Manual Deployment

### 1. Create Namespace

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: go-meta-redirector
```

```bash
kubectl apply -f k8s/namespace.yaml
```

### 2. Create ConfigMap for repos.yaml

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-meta-redirector-config
  namespace: go-meta-redirector
data:
  repos.yaml: |
    domains:
      go.pilab.hu:
        fallback:
          pattern: "cloud/*"
          target: "https://github.com/pilab-dev/*"
        repos:
          - path: cloud/log
            git_url: https://github.com/pilab-dev/log.git
            pkgsite_url: https://pkg.go.dev/go.pilab.hu/cloud/log
      go.paalgyula.com:
        fallback:
          pattern: "tools/*"
          target: "https://github.com/paalgyula/*"
        repos: []
      go.pira.hu:
        fallback: ~
        repos: []
```

```bash
kubectl apply -f k8s/configmap.yaml
```

### 3. Create Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-meta-redirector
  namespace: go-meta-redirector
  labels:
    app: go-meta-redirector
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-meta-redirector
  template:
    metadata:
      labels:
        app: go-meta-redirector
    spec:
      containers:
      - name: go-meta-redirector
        image: ghcr.io/pilab-dev/go-meta-redirector:latest
        ports:
        - containerPort: 8080
          name: http
        volumeMounts:
        - name: config
          mountPath: /etc/go-meta-redirector
          readOnly: true
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        livenessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: config
        configMap:
          name: go-meta-redirector-config
```

```bash
kubectl apply -f k8s/deployment.yaml
```

### 4. Create Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: go-meta-redirector
  namespace: go-meta-redirector
spec:
  selector:
    app: go-meta-redirector
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  type: ClusterIP
```

```bash
kubectl apply -f k8s/service.yaml
```

### 5. Create Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-meta-redirector
  namespace: go-meta-redirector
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"  # If using cert-manager
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - go.pilab.hu
    - go.paalgyula.com
    - go.pira.hu
    secretName: go-meta-redirector-tls
  rules:
  - host: go.pilab.hu
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-meta-redirector
            port:
              number: 80
  - host: go.paalgyula.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-meta-redirector
            port:
              number: 80
  - host: go.pira.hu
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-meta-redirector
            port:
              number: 80
```

```bash
kubectl apply -f k8s/ingress.yaml
```

## Using Helm (Recommended)

### Install with Helm

```bash
# Add Helm repo (if published)
helm repo add go-meta-redirector https://pilab-dev.github.io/go-meta-redirector/
helm install go-meta-redirector go-meta-redirector/go-meta-redirector

# Or use local chart
helm install go-meta-redirector ./helm/
```

### Helm Values

```yaml
# values.yaml
replicaCount: 2

image:
  repository: ghcr.io/pilab-dev/go-meta-redirector
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80
  targetPort: 8080

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: go.pilab.hu
      paths:
        - path: /
          pathType: Prefix
    - host: go.paalgyula.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: go-meta-redirector-tls
      hosts:
        - go.pilab.hu
        - go.paalgyula.com

config:
  repos.yaml: |
    domains:
      go.pilab.hu:
        fallback:
          pattern: "cloud/*"
          target: "https://github.com/pilab-dev/*"
        repos: []

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 50m
    memory: 64Mi
```

## Google Kubernetes Engine (GKE) Deployment

### 1. Create GKE Cluster

```bash
# Create cluster
gcloud container clusters create go-meta-redirector \
  --zone us-central1-a \
  --num-nodes 2 \
  --machine-type e2-small

# Get credentials
gcloud container clusters get-credentials go-meta-redirector \
  --zone us-central1-a
```

### 2. Deploy with Config

```bash
# Create configmap from local file
kubectl create configmap go-meta-redirector-config \
  --from-file=repos.yaml \
  -n go-meta-redirector

# Apply manifests
kubectl apply -f k8s/
```

### 3. Set Up Ingress with GKE

```bash
# Install nginx ingress controller (if not using GKE Ingress)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

# Or use GKE Ingress (managed)
# Update ingress.yaml to use: kubernetes.io/ingress.class: "gce"
```

### 4. Configure DNS

```bash
# Get ingress IP
kubectl get ingress go-meta-redirector -n go-meta-redirector

# Point your domains to the external IP
# go.pilab.hu    A    <EXTERNAL-IP>
# go.paalgyula.com    A    <EXTERNAL-IP>
```

## Scaling

```bash
# Scale deployment
kubectl scale deployment go-meta-redirector --replicas=5 -n go-meta-redirector

# Auto-scaling with HPA
kubectl autoscale deployment go-meta-redirector \
  --cpu-percent=80 \
  --min=2 \
  --max=10 \
  -n go-meta-redirector
```

## Monitoring

```bash
# View logs
kubectl logs -f deployment/go-meta-redirector -n go-meta-redirector

# Check resource usage
kubectl top pods -n go-meta-redirector

# Describe pod for debugging
kubectl describe pod <pod-name> -n go-meta-redirector
```

## Updating

```bash
# Update image
kubectl set image deployment/go-meta-redirector \
  go-meta-redirector=ghcr.io/pilab-dev/go-meta-redirector:v1.1.0 \
  -n go-meta-redirector

# Update config
kubectl create configmap go-meta-redirector-config \
  --from-file=repos.yaml \
  -n go-meta-redirector \
  -o yaml --dry-run | kubectl apply -f -

# Restart deployment to pick up new config
kubectl rollout restart deployment/go-meta-redirector -n go-meta-redirector
```

## Cleanup

```bash
# Delete all resources
kubectl delete namespace go-meta-redirector

# Or delete individual resources
kubectl delete -f k8s/
```

## Troubleshooting

### Pod not starting
```bash
kubectl describe pod <pod-name> -n go-meta-redirector
kubectl logs <pod-name> -n go-meta-redirector
```

### Config not loading
```bash
# Check configmap
kubectl get configmap go-meta-redirector-config -n go-meta-redirector -o yaml

# Verify mount
kubectl exec -it <pod-name> -n go-meta-redirector -- cat /etc/go-meta-redirector/repos.yaml
```

### Ingress issues
```bash
kubectl describe ingress go-meta-redirector -n go-meta-redirector
kubectl get events -n go-meta-redirector --sort-by='.lastTimestamp'
```
