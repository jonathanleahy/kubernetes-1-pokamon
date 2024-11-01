# Get minikube service URL
export SERVICE_URL=$(minikube service pokemon-service --url)

for i in {1..1000000}; do curl -s http://$SERVICE_IP > /dev/null & done

# Then run the curl loop
while true; do
    echo "$(date '+%H:%M:%S')"
    curl -s $SERVICE_URL | grep -A1 "pod-info" | grep "Pod"
    echo "---"
    sleep 1
done