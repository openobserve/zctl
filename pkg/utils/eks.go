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
	fmt.Println("..............Getting EKS Cluster Details............")
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

	fmt.Println("Got Cluster Details. Will use these details to create the IAM role for the cluster.")

	return resp.Cluster, nil
}

func HasOIDCProvider(clusterName, region string) (bool, error) {
	fmt.Println("..............Checking if an OIDC Provider already exists for the cluster............")
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
		fmt.Println("Cluster does not have OIDC provider configured. Run 'eksctl utils associate-iam-oidc-provider --region="+region+" --cluster=", clusterName, " --approve' to associate an OIDC provider with the cluster")
		return false, nil
	}
}

func GetEksClusterNameByApiServerUrl(apiServerUrl string) (string, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Create an EKS client
	svc := eks.NewFromConfig(cfg)

	// List all EKS clusters
	resp, err := svc.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		return "", err
	}

	// Loop through each cluster and compare the API server URL
	for _, clusterName := range resp.Clusters {
		clusterInfo, err := svc.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: aws.String(clusterName),
		})
		if err != nil {
			return "", err
		}
		if *clusterInfo.Cluster.Endpoint == apiServerUrl {
			return clusterName, nil
		}
	}

	return "", fmt.Errorf("could not find an EKS cluster with the API server URL %s", apiServerUrl)
}

func GetCurrentEKSClusterName() (string, error) {
	// Get the API server endpoint of the current kube context
	apiEndpoint, err := GetCurrentKubeContextAPIEndpoint()
	if err != nil {
		return "", err
	}

	apiServerUrl := apiEndpoint
	clusterName, err := GetEksClusterNameByApiServerUrl(apiServerUrl)
	if err != nil {
		return "", err
	}

	return clusterName, nil
}
