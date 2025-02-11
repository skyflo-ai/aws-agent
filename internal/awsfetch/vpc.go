package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// VPC represents a Virtual Private Cloud.
type VPC struct {
	VpcID string `json:"VpcId"`
}

// Subnet represents a VPC subnet.
type Subnet struct {
	SubnetID string `json:"SubnetId"`
	VpcID    string `json:"VpcId"`
}

// RouteTable represents a VPC route table.
type RouteTable struct {
	RouteTableID string `json:"RouteTableId"`
}

// NATGateway represents a NAT gateway.
type NATGateway struct {
	NatGatewayID string `json:"NatGatewayId"`
	SubnetID     string `json:"SubnetId"`
	VpcID        string `json:"VpcId"`
}

// InternetGateway represents an Internet gateway.
type InternetGateway struct {
	InternetGatewayID string `json:"InternetGatewayId"`
}

// FetchVPCData retrieves VPCs, Subnets, Route Tables, NAT and Internet Gateways.
func FetchVPCData(ctx context.Context, cfg aws.Config) ([]VPC, []Subnet, []RouteTable, []NATGateway, []InternetGateway, error) {
	client := ec2.NewFromConfig(cfg)

	// Fetch VPCs
	vpcOut, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error fetching VPCs: %w", err)
	}
	var vpcs []VPC
	for _, v := range vpcOut.Vpcs {
		vpcs = append(vpcs, VPC{VpcID: *v.VpcId})
	}

	// Fetch Subnets
	subnetOut, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error fetching subnets: %w", err)
	}
	var subnets []Subnet
	for _, s := range subnetOut.Subnets {
		subnet := Subnet{SubnetID: *s.SubnetId}
		if s.VpcId != nil {
			subnet.VpcID = *s.VpcId
		}
		subnets = append(subnets, subnet)
	}

	// Fetch Route Tables
	rtOut, err := client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error fetching route tables: %w", err)
	}
	var routeTables []RouteTable
	for _, rt := range rtOut.RouteTables {
		routeTables = append(routeTables, RouteTable{RouteTableID: *rt.RouteTableId})
	}

	// Fetch NAT Gateways
	natOut, err := client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error fetching NAT gateways: %w", err)
	}
	var natGateways []NATGateway
	for _, nat := range natOut.NatGateways {
		ng := NATGateway{NatGatewayID: *nat.NatGatewayId}
		if nat.SubnetId != nil {
			ng.SubnetID = *nat.SubnetId
		}
		if nat.VpcId != nil {
			ng.VpcID = *nat.VpcId
		}
		natGateways = append(natGateways, ng)
	}

	// Fetch Internet Gateways
	igwOut, err := client.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{})
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error fetching Internet gateways: %w", err)
	}
	var internetGateways []InternetGateway
	for _, igw := range igwOut.InternetGateways {
		internetGateways = append(internetGateways, InternetGateway{InternetGatewayID: *igw.InternetGatewayId})
	}

	return vpcs, subnets, routeTables, natGateways, internetGateways, nil
}
