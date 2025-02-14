# AWS Resource Crawler & Watcher

## Prerequisites

- **Go:** Version 1.23.6 is recommended
- **Docker:** To build and run container images
- **AWS Credentials:** Ensure you have AWS credentials with read-only access to your target resources
- **CloudTrail / EventBridge (for Watcher):** The default CloudTrail configuration is used. For the watcher, you can use EventBridge with API Destinations to forward events to the watcher's endpoint
- **Backend Endpoint:** A backend service (or a dummy endpoint for testing) should be available to receive the JSON data from the crawler and watcher

## Configuration

### Environment Variables

Create a `.env` file in the project root with the following content:

```dotenv
AWS_ACCESS_KEY_ID=your_test_aws_access_key
AWS_SECRET_ACCESS_KEY=your_test_aws_secret_key
AWS_REGION=us-east-1
BACKEND_ENDPOINT=http://host.docker.internal:8181/api/aws-resources
```

### Starting the Dummy Backend (for Testing)

For testing purposes, you can use the provided dummy backend that logs received data:

```bash
go run dummy_backend.go
```

The dummy backend will start on port 8181 and log all received AWS resource data to the console.

## Building & Running

### Building the Docker Images

Two separate Docker images are available: one for the initial crawler and one for the watcher.

#### Build the Crawler Image

```bash
docker build -f dockerfile.crawler -t go-aws-crawler-crawler .
```

#### Build the Watcher Image

```bash
docker build -f dockerfile.watcher -t go-aws-crawler-watcher .
```

### Running the Containers

#### Run the Initial Crawler

The initial crawler runs once and exits after sending its aggregated AWS resource data to the backend.

```bash
docker run --env-file .env go-aws-crawler-crawler
```

#### Run the Watcher

The watcher runs continuously and listens on port 8282 for incoming events.

```bash
docker run --env-file .env --add-host=host.docker.internal:host-gateway -p 8282:8282 go-aws-crawler-watcher
```

## Testing

### Testing the Watcher

Once the watcher container is running, you can simulate a CloudTrail event using curl:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "detail-type": "AWS API Call via CloudTrail",
  "source": "aws.ec2",
  "detail": {
    "eventName": "RunInstances",
    "requestParameters": {},
    "responseElements": {
      "instancesSet": {
        "items": [
          { "instanceId": "i-test123" }
        ]
      }
    }
  }
}' http://localhost:8282/events
```

To verify the event processing, check the watcher container logs:

```bash
docker logs <container_id>
```

## Deployment

### CloudFormation Stack for One-Click Watcher Deployment

To automate the configuration of EventBridge and API Destinations for the watcher, use the provided CloudFormation template (see separate documentation or file). This stack can be deployed by your client via the AWS CloudFormation console with a one-click solution.
