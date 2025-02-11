package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMUser represents an IAM user.
type IAMUser struct {
	UserName string `json:"UserName"`
	// Additional fields can be added.
}

// IAMPolicy represents an IAM policy.
type IAMPolicy struct {
	PolicyName string `json:"PolicyName"`
	// Additional fields can be added.
}

// FetchIAMData retrieves IAM users and local IAM policies.
func FetchIAMData(ctx context.Context, cfg aws.Config) ([]IAMUser, []IAMPolicy, error) {
	client := iam.NewFromConfig(cfg)
	
	// Fetch IAM Users
	userOut, err := client.ListUsers(ctx, &iam.ListUsersInput{})
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching IAM users: %w", err)
	}
	var users []IAMUser
	for _, user := range userOut.Users {
		users = append(users, IAMUser{UserName: *user.UserName})
	}

	// Fetch local IAM Policies
	policyOut, err := client.ListPolicies(ctx, &iam.ListPoliciesInput{
		Scope: "Local",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching IAM policies: %w", err)
	}
	var policies []IAMPolicy
	for _, policy := range policyOut.Policies {
		policies = append(policies, IAMPolicy{PolicyName: *policy.PolicyName})
	}

	return users, policies, nil
}
