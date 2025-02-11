package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// Route53Zone represents a Route53 hosted zone.
type Route53Zone struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

// FetchRoute53Zones retrieves all Route53 hosted zones.
func FetchRoute53Zones(ctx context.Context, cfg aws.Config) ([]Route53Zone, error) {
	client := route53.NewFromConfig(cfg)
	out, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		return nil, fmt.Errorf("error fetching Route53 zones: %w", err)
	}
	var zones []Route53Zone
	for _, zone := range out.HostedZones {
		zones = append(zones, Route53Zone{
			Id:   *zone.Id,
			Name: *zone.Name,
		})
	}
	return zones, nil
}
