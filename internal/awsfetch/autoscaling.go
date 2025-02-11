package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

// AutoScalingGroup represents an AutoScaling group.
type AutoScalingGroup struct {
	AutoScalingGroupName string `json:"AutoScalingGroupName"`
	// Add more fields if needed.
}

// FetchAutoScalingGroups retrieves all AutoScaling groups.
func FetchAutoScalingGroups(ctx context.Context, cfg aws.Config) ([]AutoScalingGroup, error) {
	client := autoscaling.NewFromConfig(cfg)
	out, err := client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching AutoScaling groups: %w", err)
	}
	var groups []AutoScalingGroup
	for _, group := range out.AutoScalingGroups {
		groups = append(groups, AutoScalingGroup{AutoScalingGroupName: *group.AutoScalingGroupName})
	}
	return groups, nil
}
