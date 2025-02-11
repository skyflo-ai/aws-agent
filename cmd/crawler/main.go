package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/DavisAndn/go-aws-crawler/internal/awsfetch"
	"github.com/DavisAndn/go-aws-crawler/internal/backend"
	"github.com/DavisAndn/go-aws-crawler/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
)

type InitialData struct {
	EC2Instances      []awsfetch.EC2Instance      `json:"ec2_instances"`
	VPCs              []awsfetch.VPC              `json:"vpcs"`
	Subnets           []awsfetch.Subnet           `json:"subnets"`
	RouteTables       []awsfetch.RouteTable       `json:"route_tables"`
	NATGateways       []awsfetch.NATGateway       `json:"nat_gateways"`
	InternetGateways  []awsfetch.InternetGateway  `json:"internet_gateways"`
	S3Buckets         []awsfetch.S3Bucket         `json:"s3_buckets"`
	RDSInstances      []awsfetch.RDSInstance      `json:"rds_instances"`
	Route53Zones      []awsfetch.Route53Zone      `json:"route53_hosted_zones"`
	AutoScalingGroups []awsfetch.AutoScalingGroup `json:"autoscaling_groups"`
	LoadBalancers     []awsfetch.LoadBalancer     `json:"load_balancers"`
	EKSClusters       []awsfetch.EKSCluster       `json:"eks_clusters"`
	IAMUsers          []awsfetch.IAMUser          `json:"iam_users"`
	IAMPolicies       []awsfetch.IAMPolicy        `json:"iam_policies"`
	ElastiCaches      []awsfetch.ElastiCache      `json:"elastic_caches"`
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Use default credential chain
	awsConfig, err := awsCfg.LoadDefaultConfig(context.Background(),
		awsCfg.WithRegion(cfg.AWSRegion),
		awsCfg.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody), // Add debug logging
	)
	if err != nil {
		log.Fatalf("Error loading AWS SDK config: %v", err)
	}

	// Add debug logging
	log.Printf("AWS Region: %s", awsConfig.Region)
	creds, err := awsConfig.Credentials.Retrieve(context.Background())
	if err != nil {
		log.Printf("Error retrieving credentials: %v", err)
	} else {
		log.Printf("AWS Credentials Provider: %s", creds.Source)
		log.Printf("Access Key ID: %s", creds.AccessKeyID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// Fetch resources concurrently
	var wg sync.WaitGroup
	initialData := &InitialData{}
	var fetchErr error
	var mu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		instances, err := awsfetch.FetchEC2Instances(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.EC2Instances = instances
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		vpcs, subnets, routeTables, natGateways, internetGateways, err := awsfetch.FetchVPCData(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.VPCs = vpcs
		initialData.Subnets = subnets
		initialData.RouteTables = routeTables
		initialData.NATGateways = natGateways
		initialData.InternetGateways = internetGateways
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		buckets, err := awsfetch.FetchS3Buckets(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.S3Buckets = buckets
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rds, err := awsfetch.FetchRDSInstances(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.RDSInstances = rds
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		zones, err := awsfetch.FetchRoute53Zones(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.Route53Zones = zones
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		groups, err := awsfetch.FetchAutoScalingGroups(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.AutoScalingGroups = groups
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lbs, err := awsfetch.FetchLoadBalancers(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.LoadBalancers = lbs
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		eks, err := awsfetch.FetchEKSClusters(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.EKSClusters = eks
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		users, policies, err := awsfetch.FetchIAMData(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.IAMUsers = users
		initialData.IAMPolicies = policies
		mu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		caches, err := awsfetch.FetchElastiCaches(ctx, awsConfig)
		if err != nil {
			fetchErr = err
			return
		}
		mu.Lock()
		initialData.ElastiCaches = caches
		mu.Unlock()
	}()

	wg.Wait()
	if fetchErr != nil {
		log.Fatalf("Error during resource fetching: %v", fetchErr)
	}

	payload, err := json.Marshal(initialData)
	if err != nil {
		log.Fatalf("Error marshaling initial data: %v", err)
	}

	err = backend.SendInitialCrawlResults(cfg.BackendEndpoint, payload)
	if err != nil {
		log.Fatalf("Error sending initial crawl results: %v", err)
	}

	log.Println("Initial crawl completed and results sent successfully.")
}
