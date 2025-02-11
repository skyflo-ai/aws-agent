#!/bin/sh
echo "Starting entrypoint script..."
echo "Current environment variables:"
env | grep AWS
echo "Setting AWS credentials..."
# Map our custom env vars to standard AWS SDK env vars
export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY

echo "Updated environment variables:"
env | grep AWS

echo "Starting crawler..."
exec /crawler 