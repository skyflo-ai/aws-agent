# AWS Resource Crawler & Watcher

A system for crawling AWS resources and monitoring changes using two CloudFormation stacks-one for crawling and one for watching CloudTrail events. This solution deploys a container-based Lambda function that gathers AWS resource data and a CloudWatch-based watcher that forwards CloudTrail events to your backend endpoint.

## Deployment and Testing Instructions

### 1. Deploying the Crawler Stack

#### Docker Image Build

Build the crawler image using:

```bash
docker build --platform linux/amd64 -f dockerfile.crawler -t go-aws-crawler-crawler .
```

**Important:**  
Use the `--platform linux/amd64` flag because macOS defaults to ARM64, which is incompatible with an x86_64 Lambda.

Alternatively, pull the pre-built image from:

```bash
public.ecr.aws/x6v5w6d9/aws-crawler:latest
```

#### CloudFormation Crawler Stack (`crawler-stack.yaml`)

Deploying this stack creates the following resources:

- **Private ECR Repository:**  
  A repository named `skyflo-aws-crawler-1` to store the crawler image.
- **CodeBuild Project:**  
  Pulls the public image, tags it as `skyflo-aws-crawler-1:latest`, and pushes it to the private ECR.
- **Inline Polling Lambda Function:**  
  Triggers the CodeBuild project and polls until the build completes, sending a response back to CloudFormation.
- **Crawler Lambda Function:**  
  Once the image is available in the private ECR, a container-based Lambda function is created with a Function URL for direct invocation.

#### Triggering the Crawler

After deployment, you can trigger the crawler Lambda either via its Function URL or through the Lambda console. When executed, the crawler gathers AWS resource data and posts a JSON payload with the results to your backend endpoint.

### 2. Deploying the Watcher Stack

Deploy the `watcher-stack.yaml` using AWS CloudFormation. This stack creates:

- **EventBridge Connection & API Destination:**  
  Establishes a secure connection (using API_KEY authentication with a dummy key by default) to your backend endpoint.
- **EventBridge Rule:**  
  Listens for CloudTrail events from multiple AWS services (e.g., EC2, VPC, IAM, AutoScaling, ELB, EKS, ElastiCache, Route53, S3) and forwards them to your backend.

### 3. Testing the Setup

#### Backend Setup for Testing

For local testing, run the dummy backend:

```bash
go run dummy_backend.go
```

Use [ngrok](https://ngrok.com/) to expose the local backend (running on port 8181) to the internet, and update the backend URL in both YAML files accordingly.

#### Testing the Crawler

- Trigger the crawler Lambda function via its Function URL or the Lambda console.
- Verify that the dummy backend logs display the received JSON payload containing AWS resource data.

#### Testing the Watcher

- Simulate changes in your AWS environment or generate CloudTrail events.
- Confirm that these events are forwarded to your backend endpoint by checking the dummy backend logs.

## Summary

- **Crawler:**

  - **Build Process:** Utilizes CodeBuild and an inline polling Lambda to copy a Docker image from a public repository into a private ECR repository.
  - **Deployment:** A container-based Lambda function with a Function URL is deployed.
  - **Execution:** When triggered, it gathers AWS resource data and sends the results to the backend.

- **Watcher:**
  - **Monitoring:** Uses an EventBridge rule to capture CloudTrail events from various AWS services.
  - **Notification:** Forwards any detected changes to the backend via an API Destination.

Ensure that your backend endpoint (e.g., your ngrok URL) is correctly set in both CloudFormation YAML files before deployment.
