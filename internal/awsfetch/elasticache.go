package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

// ElastiCache represents an ElastiCache cluster.
type ElastiCache struct {
	CacheClusterId string `json:"CacheClusterId"`
	// Additional fields can be added.
}

// FetchElastiCaches retrieves all ElastiCache clusters.
func FetchElastiCaches(ctx context.Context, cfg aws.Config) ([]ElastiCache, error) {
	client := elasticache.NewFromConfig(cfg)
	out, err := client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		ShowCacheNodeInfo: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching ElastiCache clusters: %w", err)
	}
	var caches []ElastiCache
	for _, cache := range out.CacheClusters {
		caches = append(caches, ElastiCache{CacheClusterId: *cache.CacheClusterId})
	}
	return caches, nil
}
