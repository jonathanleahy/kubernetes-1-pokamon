#!/bin/bash

# Set your Docker Hub username
DOCKER_USERNAME="jonathanleahy"
APP_NAME="pokemon-checker"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}>>> $1${NC}"
}

# Function to print error messages
print_error() {
    echo -e "${RED}ERROR: $1${NC}"
    exit 1
}

# Check if Docker Hub username is set
if [ "$DOCKER_USERNAME" = "YOUR_DOCKERHUB_USERNAME" ]; then
    print_error "Please set your Docker Hub username in the script"
fi

# Build Go application
print_status "Building Go application..."
go build -o main . || print_error "Go build failed"

# Build Docker image
print_status "Building Docker image..."
docker build -t $APP_NAME:latest . || print_error "Docker build failed"

# Tag Docker image
print_status "Tagging Docker image..."
docker tag $APP_NAME:latest $DOCKER_USERNAME/$APP_NAME:latest || print_error "Docker tag failed"

# Push to Docker Hub
print_status "Pushing to Docker Hub..."
docker push $DOCKER_USERNAME/$APP_NAME:latest || print_error "Docker push failed"

# Apply Kubernetes configurations
print_status "Applying Kubernetes configurations..."

# Apply deployment
kubectl apply -f deployment.yaml || print_error "Failed to apply deployment"

# Apply service
kubectl apply -f service.yaml || print_error "Failed to apply service"

# Restart deployment to pull new image
print_status "Restarting deployment to pull new image..."
kubectl rollout restart deployment $APP_NAME || print_error "Failed to restart deployment"

# Wait for rollout to complete
print_status "Waiting for rollout to complete..."
kubectl rollout status deployment/$APP_NAME || print_error "Rollout failed"

# Get service information
print_status "Getting service information..."
kubectl get service $APP_NAME-service

print_status "Deployment completed successfully!"

# Display pod status
echo ""
print_status "Pod status:"
kubectl get pods -l app=$APP_NAME

echo ""
print_status "You can watch pod status with: kubectl get pods -l app=$APP_NAME -w"
print_status "Access the application using the External-IP from the service output above"
