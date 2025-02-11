package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

// LoadBalancer represents a load balancer (both classic and modern).
type LoadBalancer struct {
	LoadBalancerName string `json:"LoadBalancerName"`
	// Add additional fields (e.g., VpcId, Type) as needed.
}

// FetchLoadBalancers retrieves load balancers from both ELB and ELBv2.
func FetchLoadBalancers(ctx context.Context, cfg aws.Config) ([]LoadBalancer, error) {
	var lbs []LoadBalancer

	// Modern load balancers using ELBv2
	clientV2 := elasticloadbalancingv2.NewFromConfig(cfg)
	outV2, err := clientV2.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching modern load balancers: %w", err)
	}
	for _, lb := range outV2.LoadBalancers {
		lbs = append(lbs, LoadBalancer{LoadBalancerName: *lb.LoadBalancerName})
	}

	// Classic load balancers using ELB
	clientClassic := elasticloadbalancing.NewFromConfig(cfg)
	outClassic, err := clientClassic.DescribeLoadBalancers(ctx, &elasticloadbalancing.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching classic load balancers: %w", err)
	}
	for _, lb := range outClassic.LoadBalancerDescriptions {
		lbs = append(lbs, LoadBalancer{LoadBalancerName: *lb.LoadBalancerName})
	}

	return lbs, nil
}
