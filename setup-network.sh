#!/bin/bash
# Script to create the external Docker network for Studoto

NETWORK_NAME="studoto"
SUBNET="172.28.5.0/24"

echo "Creating Docker network: $NETWORK_NAME with subnet $SUBNET"

# Check if network already exists
if docker network ls | grep -q "$NETWORK_NAME"; then
    echo "Network '$NETWORK_NAME' already exists."
    read -p "Do you want to remove and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Removing existing network..."
        docker network rm "$NETWORK_NAME" 2>/dev/null || true
    else
        echo "Using existing network."
        exit 0
    fi
fi

# Create the network
docker network create --driver bridge --subnet="$SUBNET" "$NETWORK_NAME"

if [ $? -eq 0 ]; then
    echo "✅ Network '$NETWORK_NAME' created successfully!"
    echo "Subnet: $SUBNET"
else
    echo "❌ Failed to create network"
    exit 1
fi
