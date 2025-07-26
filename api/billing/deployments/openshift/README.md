# OpenShift Deployment Guide

This directory contains OpenShift deployment manifests for the billing service, now configured to use GitHub Container Registry instead of local image builds.

## Quick Start

### Prerequisites

1. **OpenShift CLI** installed and configured
2. **Access to your OpenShift cluster**
3. **GitHub Personal Access Token** with package read permissions
4. **Container images** published to GitHub Container Registry via CI/CD

### Setup GitHub Container Registry Access

Create a secret for pulling images from GitHub Container Registry:

```bash
# Create pull secret for GitHub Container Registry
oc create secret docker-registry ghcr-pull-secret \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password=YOUR_GITHUB_TOKEN \
  --docker-email=YOUR_EMAIL \
  -n billing-staging

oc create secret docker-registry ghcr-pull-secret \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password=YOUR_GITHUB_TOKEN \
  --docker-email=YOUR_EMAIL \
  -n billing-production

# Link the secret to the service accounts
oc secrets link billing-api ghcr-pull-secret --for=pull -n billing-staging
oc secrets link billing-migrations ghcr-pull-secret --for=pull -n billing-staging
oc secrets link billing-api ghcr-pull-secret --for=pull -n billing-production
oc secrets link billing-migrations ghcr-pull-secret --for=pull -n billing-production
```

### Deploy to Staging

```bash
# Create staging namespace
oc apply -f staging/namespace.yaml

# Deploy core services (if not already deployed)
oc apply -f secret.yaml -n billing-staging
oc apply -f configmap.yaml -n billing-staging
oc apply -f service.yaml -n billing-staging

# Deploy application
oc apply -f staging/deployment.yaml

# Create route
oc apply -f route.yaml -n billing-staging

# Verify deployment
oc get pods -n billing-staging
oc get routes -n billing-staging
```

### Deploy to Production

```bash
# Create production namespace
oc apply -f production/namespace.yaml

# Deploy core services (if not already deployed)
oc apply -f secret.yaml -n billing-production
oc apply -f configmap.yaml -n billing-production
oc apply -f service.yaml -n billing-production

# Deploy application
oc apply -f production/deployment.yaml

# Create route
oc apply -f route.yaml -n billing-production

# Verify deployment
oc get pods -n billing-production
oc get routes -n billing-production
```

## Deployment Order

For manual deployment, apply manifests in this order:

```bash
# 1. Create namespace
oc apply -f namespace.yaml

# 2. Create secrets and config
oc apply -f secret.yaml
oc apply -f configmap.yaml

# 3. Create image streams
oc apply -f imagestream.yaml

# 4. Create build configs (will trigger builds)
oc apply -f buildconfig.yaml

# 5. Wait for builds to complete
oc logs -f bc/billing-api -n billing
oc logs -f bc/billing-migrations -n billing

# 6. Deploy the application
oc apply -f deploymentconfig.yaml
oc apply -f service.yaml
oc apply -f route.yaml
```

## Configuration

### Database Connection

Update the secret with your PostgreSQL credentials:

```bash
# Edit the secret
oc edit secret billing-secrets -n billing

# Or create from command line
oc create secret generic billing-secrets \
  --from-literal=BILLING_DATABASE_PASSWORD="your-password" \
  --from-literal=BILLING_MIGRATE_DATABASE_PASSWORD="your-password" \
  -n billing
```

### Git Repository

Update `buildconfig.yaml` with your Git repository URL:

```yaml
source:
  git:
    uri: https://github.com/your-org/gotuto.git
    ref: main
```

### Custom Domain

Update `route.yaml` if you want a custom hostname:

```yaml
spec:
  host: billing-api.your-domain.com
```

## Build Process

The build process uses OpenShift's Source-to-Image (S2I) strategy with Red Hat UBI base images:

1. **billing-api**: Builds the Go API service
2. **billing-migrations**: Builds the database migration tool

### Triggering Builds

```bash
# Manual build trigger
oc start-build billing-api -n billing

# Follow build logs
oc logs -f bc/billing-api -n billing

# Check build status
oc get builds -n billing
```

## Monitoring and Debugging

### Check Pod Status
```bash
oc get pods -n billing
oc describe pod <pod-name> -n billing
```

### View Logs
```bash
# API logs
oc logs -f deployment/billing-api -n billing

# Migration logs (from init container)
oc logs <pod-name> -c migrate -n billing
```

### Debug Networking
```bash
# Test service connectivity
oc exec -it <pod-name> -n billing -- curl http://billing-api-service:80/health

# Check route
curl https://$(oc get route billing-api -n billing -o jsonpath='{.spec.host}')/health
```

## Security Features

- **Non-root containers**: All containers run as non-root user (UID 1001)
- **Security contexts**: Proper security contexts for OpenShift
- **Secrets management**: Database credentials stored as secrets
- **TLS termination**: HTTPS enabled via OpenShift routes
- **Network policies**: (can be added for micro-segmentation)

## Scaling

```bash
# Scale the API service
oc scale dc billing-api --replicas=5 -n billing

# Check scaling status
oc get pods -n billing -l app=billing-api
```

## Rolling Updates

OpenShift automatically handles rolling updates when you push new code:

```bash
# Force a new deployment
oc rollout latest dc/billing-api -n billing

# Check rollout status
oc rollout status dc/billing-api -n billing

# Rollback if needed
oc rollout undo dc/billing-api -n billing
```

## Resource Limits

The deployment includes resource requests and limits:

- **Requests**: 128Mi memory, 100m CPU
- **Limits**: 256Mi memory, 200m CPU

Adjust these in `deploymentconfig.yaml` based on your needs.

## Health Checks

The deployment includes both liveness and readiness probes:

- **Liveness**: `/health` endpoint (restarts pod if unhealthy)
- **Readiness**: `/health` endpoint (removes from load balancer if not ready)

## Migration Strategy

Database migrations run as init containers:

1. Migration init container runs first
2. Applies all pending migrations
3. API container starts only after migrations succeed
4. If migrations fail, pod doesn't start

## Environment-specific Configuration

For different environments (dev, staging, prod):

1. Create separate namespaces
2. Use different ConfigMaps for environment-specific settings
3. Use different Secrets for environment-specific credentials
4. Adjust resource limits per environment

## Troubleshooting

### Build Failures
```bash
# Check build logs
oc logs bc/billing-api -n billing

# Common issues:
# - Git repository access
# - Dockerfile/Containerfile syntax
# - Missing dependencies
```

### Deployment Failures
```bash
# Check events
oc get events -n billing --sort-by='.lastTimestamp'

# Check pod description
oc describe pod <failing-pod> -n billing

# Common issues:
# - Image pull failures
# - Secret/ConfigMap missing
# - Resource constraints
```

### Database Connection Issues
```bash
# Test from pod
oc exec -it <pod-name> -n billing -- env | grep BILLING_DATABASE

# Check if PostgreSQL service is available
oc get svc -n billing
```

This OpenShift deployment provides a production-ready setup for your billing service with proper security, monitoring, and scalability features.