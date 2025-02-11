package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// RDSInstance represents an RDS instance.
type RDSInstance struct {
	DBInstanceIdentifier string `json:"DBInstanceIdentifier"`
	// Add additional fields as needed.
}

// FetchRDSInstances retrieves all RDS instances.
func FetchRDSInstances(ctx context.Context, cfg aws.Config) ([]RDSInstance, error) {
	client := rds.NewFromConfig(cfg)
	out, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching RDS instances: %w", err)
	}
	var instances []RDSInstance
	for _, db := range out.DBInstances {
		instances = append(instances, RDSInstance{DBInstanceIdentifier: *db.DBInstanceIdentifier})
	}
	return instances, nil
}
