package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2Instance represents a simplified EC2 instance.
type EC2Instance struct {
	InstanceID string            `json:"InstanceId"`
	VpcID      string            `json:"VpcId"`
	SubnetID   string            `json:"SubnetId"`
	Tags       map[string]string `json:"Tags"`
}

// FetchEC2Instances retrieves all EC2 instances.
func FetchEC2Instances(ctx context.Context, cfg aws.Config) ([]EC2Instance, error) {
	client := ec2.NewFromConfig(cfg)
	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{})
	var instances []EC2Instance

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching EC2 instances: %w", err)
		}
		for _, reservation := range page.Reservations {
			for _, inst := range reservation.Instances {
				instance := EC2Instance{
					InstanceID: *inst.InstanceId,
				}
				if inst.VpcId != nil {
					instance.VpcID = *inst.VpcId
				}
				if inst.SubnetId != nil {
					instance.SubnetID = *inst.SubnetId
				}
				tagsMap := make(map[string]string)
				for _, tag := range inst.Tags {
					if tag.Key != nil && tag.Value != nil {
						tagsMap[*tag.Key] = *tag.Value
					}
				}
				instance.Tags = tagsMap
				instances = append(instances, instance)
			}
		}
	}
	return instances, nil
}
