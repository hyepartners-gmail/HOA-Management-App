#!/bin/bash
set -e

# ---- CONFIG ----
PROJECT_ID="hyepartners-324923474516"
REGION="us-central1"
REPO="hoa-management-app"
SERVICE_NAME="hoamanager"
IMAGE_TAG="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/${SERVICE_NAME}"

# ---- BUILD + TAG LOCALLY ----
echo "🔨 Building Docker image..."
docker build -t "$IMAGE_TAG" .

# ---- AUTH TO GCP CONTAINER REGISTRY ----
echo "🔐 Authenticating Docker to GCP..."
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

# ---- PUSH TO ARTIFACT REGISTRY ----
echo "🚀 Pushing to Artifact Registry..."
docker push "$IMAGE_TAG"

# ---- DEPLOY TO CLOUD RUN ----
echo "🌐 Deploying to Cloud Run..."
gcloud run deploy "$SERVICE_NAME" \
  --image="$IMAGE_TAG" \
  --platform=managed \
  --region="$REGION" \
  --allow-unauthenticated

echo "✅ Deployment complete."
