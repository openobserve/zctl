package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func GetEKSClusterDetails(clusterName string) (*types.Cluster, error) {
	fmt.Println("..............GetEKSClusterDetails............")
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Create a new EKS client
	svc := eks.NewFromConfig(cfg)

	// Call the DescribeCluster API to retrieve the cluster details
	resp, err := svc.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return nil, err
	}

	return resp.Cluster, nil
}

func HasOIDCProvider(clusterName string) (bool, error) {
	fmt.Println("..............HasOIDCProvider............")
	// Load the AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return false, err
	}

	// Create a new Amazon EKS client
	svc := eks.NewFromConfig(cfg)

	// Call the DescribeCluster API to get information about the cluster
	resp, err := svc.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})
	if err != nil {
		return false, err
	}

	// Check if the cluster has an OIDC provider configured
	if resp.Cluster.Identity.Oidc != nil {
		fmt.Println("Cluster has OIDC provider configured")
		return true, nil
	} else {
		return false, nil
	}
}

// func CheckOIDCProviderExistsForCluster(clusterName string) (bool, error) {
// 	fmt.Println("..............CheckOIDCProviderExistsForCluster............")
// 	// Load the AWS configuration
// 	cfg, err := config.LoadDefaultConfig(context.Background())
// 	if err != nil {
// 		return false, err
// 	}

// 	// Create a new EKS client
// 	svc := eks.NewFromConfig(cfg)

// 	// idpConfig := &types.IdentityProviderConfig{
// 	// 	Name: aws.String("oidc"),
// 	// 	Type: aws.String("oidc"),
// 	// }

// 	res, err := svc.ListIdentityProviderConfigs(context.Background(), &eks.ListIdentityProviderConfigsInput{
// 		ClusterName: aws.String(clusterName),
// 	})

// 	if err != nil {
// 		return false, err
// 	}

// 	ip, _ := json.Marshal(res)

// 	fmt.Println("res: ", string(ip))

// 	for _, idpConfig := range res.IdentityProviderConfigs {
// 		fmt.Println("Name: ", idpConfig.Name)
// 		if *idpConfig.Name == "oidc" {
// 			fmt.Println("Found OIDC provider", idpConfig.Name)
// 			return true, nil
// 		}
// 	}

// 	// Call the DescribeIdentityProviderConfig API to retrieve the OIDC identity provider configuration
// 	// resp, err := svc.DescribeIdentityProviderConfig(context.Background(), &eks.DescribeIdentityProviderConfigInput{
// 	// 	ClusterName: aws.String(clusterName),
// 	// 	IdentityProviderConfig: &types.IdentityProviderConfig{
// 	// 		Name: aws.String("oidc"),
// 	// 		Type: aws.String("oidc"),
// 	// 	},
// 	// })
// 	// if err != nil {
// 	// 	return false, err
// 	// }

// 	// // Check if an OIDC identity provider is configured for the cluster
// 	// if resp.IdentityProviderConfig != nil && resp.IdentityProviderConfig.Oidc != nil {
// 	// 	return true, nil
// 	// }

// 	return false, nil
// }
