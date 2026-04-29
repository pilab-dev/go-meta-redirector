# Deploying GoMeta Redirector on Google Cloud Run

This guide covers deploying GoMeta Redirector to Google Cloud Run for a fully managed, serverless experience.

## Prerequisites

- Google Cloud account
- `gcloud` CLI installed and authenticated
- Docker image pushed to `ghcr.io/pilab-dev/go-meta-redirector:latest`
- Domain names ready to configure

## Quick Deploy

```bash
# Set project
gcloud config set project YOUR_PROJECT_ID

# Deploy to Cloud Run
gcloud run deploy go-meta-redirector \
  --image ghcr.io/pilab-dev/go-meta-redirector:latest \
  --platform managed \
  --region us-central1 \
  --port 8080 \
  --allow-unauthenticated \
  --set-env-vars="PORT=8080"
```

## Detailed Deployment Steps

### 1. Prepare Cloud Run Configuration

Create `cloud-run-service.yaml`:

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: go-meta-redirector
  annotations:
    run.googleapis.com/launch-stage: GA
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/maxScale: "10"
        autoscaling.knative.dev/minScale: "1"
        run.googleapis.com/cpu-throttling: "false"
    spec:
      containers:
      - image: ghcr.io/pilab-dev/go-meta-redirector:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        resources:
          limits:
            cpu: "1"
            memory: "256Mi"
        # Mount config from Secret Manager
        volumeMounts:
        - name: config
          mountPath: /etc/go-meta-redirector
      volumes:
      - name: config
        secret:
          secretName: go-meta-redirector-config
          items:
          - key: repos.yaml
            path: repos.yaml
```

### 2. Store Configuration in Secret Manager

```bash
# Create secret with repos.yaml content
gcloud secrets create go-meta-redirector-config \
  --data-file=repos.yaml \
  --replication-policy="automatic"

# Or create from stdin
cat repos.yaml | gcloud secrets create go-meta-redirector-config \
  --data-file=- \
  --replication-policy="automatic"

# Grant Cloud Run service account access
gcloud secrets add-iam-policy-binding go-meta-redirector-config \
  --member=serviceAccount:$(gcloud projects describe YOUR_PROJECT_ID --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
  --role=roles/secretmanager.secretAccessor
```

### 3. Deploy the Service

```bash
# Deploy with gcloud
gcloud run deploy go-meta-redirector \
  --image ghcr.io/pilab-dev/go-meta-redirector:latest \
  --platform managed \
  --region us-central1 \
  --port 8080 \
  --allow-unauthenticated \
  --set-secrets="/etc/go-meta-redirector/repos.yaml=go-meta-redirector-config:latest" \
  --cpu 1 \
  --memory 256Mi \
  --min-instances 1 \
  --max-instances 10 \
  --set-env-vars="PORT=8080" \
  --timeout=300s \
  --concurrency=80
```

### 4. Configure Custom Domains

Cloud Run supports custom domains with automatic TLS:

```bash
# Map domain to service
gcloud run domain-mappings create \
  --service go-meta-redirector \
  --domain go.pilab.hu \
  --region us-central1

gcloud run domain-mappings create \
  --service go-meta-redirector \
  --domain go.paalgyula.com \
  --region us-central1

# Get required DNS records
gcloud run domain-mappings describe \
  --domain go.pilab.hu \
  --region us-central1
```

Update your DNS:
```
go.pilab.hu     CNAME   gvrhi7zp5a-uk.a.run.app
go.paalgyula.com CNAME   gvrhi7zp5a-uk.a.run.app
```

### 5. Verify Deployment

```bash
# Get service URL
gcloud run services describe go-meta-redirector \
  --region us-central1 \
  --format="value(status.url)"

# Test the service
curl "https://go.pilab.hu/cloud/log?go-get=1"

# Test with curl
curl -H "Host: go.pilab.hu" \
  "https://YOUR_SERVICE_URL/cloud/log?go-get=1"
```

## Using Cloud Build (CI/CD)

Create `cloudbuild.yaml` for automated deployments:

```yaml
steps:
  # Build the image
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/go-meta-redirector:$SHORT_SHA', '.']

  # Push to Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/go-meta-redirector:$SHORT_SHA']

  # Deploy to Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    entrypoint: gcloud
    args:
      - 'run'
      - 'deploy'
      - 'go-meta-redirector'
      - '--image'
      - 'gcr.io/$PROJECT_ID/go-meta-redirector:$SHORT_SHA'
      - '--region'
      - 'us-central1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'

images:
  - 'gcr.io/$PROJECT_ID/go-meta-redirector:$SHORT_SHA'
```

Trigger deployment:

```bash
# Submit build
gcloud builds submit --config cloudbuild.yaml

# Or connect GitHub repo for automatic builds
gcloud builds triggers create github \
  --repo-name=go-meta-redirector \
  --repo-owner=pilab-dev \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## Multiple Services for Multiple Domains

Since Cloud Run maps one domain per service, create separate services:

```bash
# Deploy for go.pilab.hu
gcloud run deploy go-meta-redirector-pilab \
  --image ghcr.io/pilab-dev/go-meta-redirector:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated

gcloud run domain-mappings create \
  --service go-meta-redirector-pilab \
  --domain go.pilab.hu \
  --region us-central1

# Deploy for go.paalgyula.com
gcloud run deploy go-meta-redirector-paalgyula \
  --image ghcr.io/pilab-dev/go-meta-redirector:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated

gcloud run domain-mappings create \
  --service go-meta-redirector-paalgyula \
  --domain go.paalgyula.com \
  --region us-central1
```

## Configuration Updates

```bash
# Update configuration secret
gcloud secrets versions add go-meta-redirector-config \
  --data-file=repos.yaml

# Redeploy to pick up new config
gcloud run services update go-meta-redirector \
  --region us-central1 \
  --set-secrets="/etc/go-meta-redirector/repos.yaml=go-meta-redirector-config:latest"
```

## Monitoring and Logging

```bash
# View logs
gcloud run services logs read go-meta-redirector \
  --region us-central1 \
  --limit=50

# Monitor in Cloud Console
# https://console.cloud.google.com/run/detail/us-central1/go-meta-redirector/logs

# Set up alerts
gcloud alpha monitoring policies create \
  --policy-from-file=alert-policy.yaml
```

## Cost Optimization

```bash
# Minimize instances (scale to zero)
gcloud run services update go-meta-redirector \
  --region us-central1 \
  --min-instances=0 \
  --max-instances=3

# CPU allocation only during request processing
gcloud run services update go-meta-redirector \
  --region us-central1 \
  --cpu-throttling

# Use smaller instance size
gcloud run services update go-meta-redirector \
  --region us-central1 \
  --memory=128Mi \
  --cpu=0.5
```

## Security

```bash
# Restrict to authenticated users only
gcloud run services update go-meta-redirector \
  --region us-central1 \
  --no-allow-unauthenticated

# Add IAM binding for public access
gcloud run services add-iam-policy-binding go-meta-redirector \
  --region us-central1 \
  --member=allUsers \
  --role=roles/run.invoker

# Use Binary Authorization (for image signature verification)
gcloud run deploy go-meta-redirector \
  --region us-central1 \
  --binary-authorization-policy=projects/$PROJECT_ID/policy/binary-auth-policy
```

## Cleanup

```bash
# Delete the service
gcloud run services delete go-meta-redirector \
  --region us-central1

# Delete domain mappings
gcloud run domain-mappings delete \
  --domain go.pilab.hu \
  --region us-central1

# Delete secret
gcloud secrets delete go-meta-redirector-config
```

## Troubleshooting

### Service not starting
```bash
# Check service status
gcloud run services describe go-meta-redirector --region us-central1

# View revision details
gcloud run revisions list --service go-meta-redirector --region us-central1
```

### Config not loading
```bash
# SSH into running container (debug)
gcloud run services proxy go-meta-redirector --region us-central1 &
curl http://localhost:8080/cloud/log?go-get=1
```

### Domain mapping issues
```bash
# Check domain mapping status
gcloud run domain-mappings list --region us-central1

# Verify DNS
dig go.pilab.hu
nslookup go.pilab.hu
```

## Example: Complete Deployment Script

```bash
#!/bin/bash
# deploy-cloud-run.sh

set -e

PROJECT_ID="your-project-id"
REGION="us-central1"
SERVICE_NAME="go-meta-redirector"
IMAGE="ghcr.io/pilab-dev/go-meta-redirector:latest"

echo "🚀 Deploying GoMeta Redirector to Cloud Run..."

# Set project
gcloud config set project $PROJECT_ID

# Deploy service
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE \
  --platform managed \
  --region $REGION \
  --port 8080 \
  --allow-unauthenticated \
  --min-instances 1 \
  --max-instances 5 \
  --cpu 1 \
  --memory 256Mi \
  --set-env-vars="PORT=8080" \
  --timeout=300s

# Get service URL
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME \
  --region $REGION \
  --format="value(status.url)")

echo "✅ Service deployed at: $SERVICE_URL"

# Map custom domains
echo "📡 Mapping custom domains..."
gcloud run domain-mappings create \
  --service $SERVICE_NAME \
  --domain go.pilab.hu \
  --region $REGION

gcloud run domain-mappings create \
  --service $SERVICE_NAME \
  --domain go.paalgyula.com \
  --region $REGION

echo "✅ Deployment complete!"
echo "Don't forget to update your DNS records."
```

Make it executable: `chmod +x deploy-cloud-run.sh` and run: `./deploy-cloud-run.sh`
