# go-aws-crawler
`docker build -t go-aws-crawler .`


`docker run --env-file .env -p 8181:8181 go-aws-crawler`

- The container will start, perform the initial crawl (fetching all AWS resources), and send the aggregated JSON payload to the backend endpoint.
- The real‑time server, it will listen on port 8181 for incoming AWS event notifications.