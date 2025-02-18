# AWS Resource Crawler & Watcher

A system for crawling AWS resources and monitoring for changes. This solution consists of two main CloudFormation stacks:

1. **Crawler Stack**  
   Gathers a comprehensive inventory of AWS resources using a container-based Lambda function.
2. **Watcher Stack**  
   Monitors AWS CloudTrail events and forwards detected changes to a backend endpoint.

---

## Architecture Overview

### Crawler Workflow

1. **Docker Image Build**  
   The crawler image is built using the following command:

   ```bash
   docker build --platform linux/amd64 -f dockerfile.crawler -t go-aws-crawler-crawler .
   ```

   **Important:** The `--platform linux/amd64` flag is required because macOS defaults to ARM64, which will not build correctly for a x86_64 Lambda.

2. **Image Source**  
   The image is already available in a public repository:

   ```bash
   public.ecr.aws/x6v5w6d9/aws-crawler:latest
   ```

   You can pull this image if you prefer not to build it locally.

3. **CloudFormation Crawler Stack**  
   When you deploy the `crawler-stack.yaml`, the following resources are created:

   - **Private ECR Repository:**  
     A repository named `skyflo-aws-crawler-1` to store the crawler image.
   - **CodeBuild Project:**  
     Pulls the public image, tags it as `skyflo-aws-crawler-1:latest`, and pushes it to the private ECR.
   - **Inline Polling Lambda Function:**  
     This inline Lambda function triggers the CodeBuild project and polls until the build completes, then sends a response back to CloudFormation.
   - **Final Crawler Lambda Function:**  
     Once the image is successfully copied to the private ECR, a container-based Lambda function is created with a Function URL for direct invocation.

4. **Triggering the Crawler**  
   After successful deployment:
   - You can trigger the crawler Lambda either via its Function URL or from the Lambda console.
   - The crawler aggregates AWS resource data and sends a JSON payload with the results to your backend endpoint.

---

### Watcher Workflow

1. **CloudFormation Watcher Stack**  
   When you deploy the `watcher-stack.yaml`, the following resources are created:

   - **EventBridge Connection & API Destination:**  
     Establishes a secure connection (using API_KEY authentication with a dummy key by default) to your backend endpoint.
   - **EventBridge Rule:**  
     Monitors CloudTrail events from multiple AWS services (EC2, VPC, IAM, AutoScaling, ELB, EKS, ElastiCache, Route53, S3, etc.) and forwards them to your backend via the API Destination.

2. **Event Flow**
   - AWS CloudTrail records events from your AWS environment.
   - The EventBridge rule picks up these events and forwards them to your backend endpoint, ensuring near real-time notification of any changes.

---

## Deployment Instructions

### 1. Deploying the Crawler Stack

1. **(Optional) Build the Docker Image Locally:**  
   If needed, build the image with:

   ```bash
   docker build --platform linux/amd64 -f dockerfile.crawler -t go-aws-crawler-crawler .
   ```

   Alternatively, pull the pre-built image from:

   ```bash
   public.ecr.aws/x6v5w6d9/aws-crawler:latest
   ```

2. **Deploy the CloudFormation Stack:**  
   Use the AWS CloudFormation console or CLI to deploy `crawler-stack.yaml`. This stack will create:

   - A private ECR repository.
   - A CodeBuild project that copies the public image into your private ECR.
   - An inline polling Lambda function to monitor the CodeBuild job.
   - The final container-based crawler Lambda function with its Function URL.

3. **Trigger the Crawler:**  
   Once the stack is deployed, trigger the crawler Lambda via its Function URL or through the Lambda console. The crawler will:
   - Aggregate AWS resource data.
   - Post the JSON results to your backend endpoint.

### 2. Deploying the Watcher Stack

1. **Deploy the CloudFormation Stack:**  
   Deploy `watcher-stack.yaml` using AWS CloudFormation. This stack creates:

   - An EventBridge connection and API Destination for secure communication with your backend.
   - An EventBridge rule that listens to CloudTrail events from multiple AWS services.

2. **Event Forwarding:**  
   Once deployed, any CloudTrail events (e.g., changes in EC2, VPC, IAM, etc.) are automatically forwarded to your backend endpoint.

---

## Testing the Setup

1. **Backend Setup for Testing:**  
   For local testing, run the dummy backend with:

   ```bash
   go run dummy_backend.go
   ```

   Use [ngrok](https://ngrok.com/) to expose the local backend (running on port 8181) to the internet.  
   Update the backend URL in both the crawler and watcher YAML files to point to your ngrok URL.

2. **Trigger and Monitor:**
   - **Crawler:**  
     Trigger the crawler Lambda function via its URL or from the console, and verify that the dummy backend logs show the received JSON payload.
   - **Watcher:**  
     Generate or simulate changes in your AWS environment to produce CloudTrail events, and confirm that these events are forwarded to your backend endpoint.

---

## Summary

- **Crawler:**

  - **Build Process:** Uses CodeBuild and an inline polling Lambda to copy a Docker image from a public repository into a private ECR repository.
  - **Deployment:** Once the image is available, a container-based Lambda function is deployed with a Function URL for on-demand invocation.
  - **Execution:** When triggered, the crawler gathers AWS resource data and sends it to the backend.

- **Watcher:**
  - **Monitoring:** Sets up an EventBridge rule to listen to CloudTrail events from various AWS services.
  - **Notification:** Forwards any detected changes to the backend via an API Destination.

Ensure that the backend endpoint (e.g., your ngrok URL) is correctly configured in both the crawler and watcher CloudFormation YAML files before deployment.
