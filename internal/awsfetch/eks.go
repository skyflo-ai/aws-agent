package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// EKSCluster represents an EKS cluster.
type EKSCluster struct {
	Name string `json:"name"`
	// Additional fields such as version, status, and resourcesVpcConfig can be added.
}

// FetchEKSClusters retrieves all EKS clusters.
func FetchEKSClusters(ctx context.Context, cfg aws.Config) ([]EKSCluster, error) {
	client := eks.NewFromConfig(cfg)
	listOut, err := client.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("error listing EKS clusters: %w", err)
	}
	var clusters []EKSCluster
	for _, name := range listOut.Clusters {
		desc, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: &name,
		})
		if err != nil {
			return nil, fmt.Errorf("error describing EKS cluster %s: %w", name, err)
		}
		clusters = append(clusters, EKSCluster{Name: *desc.Cluster.Name})
	}
	return clusters, nil
}
