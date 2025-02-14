#!/bin/sh
echo "Starting entrypoint script..."
echo "Setting AWS credentials..."
# Map our custom env vars to standard AWS SDK env vars
echo "Updated environment variables"
echo "Starting crawler..."
exec /crawler 