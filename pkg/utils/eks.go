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

// HasOIDCProvider is a function that checks whether an OIDC provider has already been configured for an Amazon EKS cluster.
// The function takes in the name of the cluster and its region as arguments.
// The function returns a boolean value indicating whether an OIDC provider exists for the cluster and an error if one occurs.
func HasOIDCProvider(clusterName, region string) (bool, error) {
	fmt.Println("..............Checking if an OIDC Provider already exists for the cluster............")
	// Load the default AWS SDK configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		// Return an error if an error occurs while loading the configuration.
		return false, err
	}

	// Create a new Amazon EKS client using the loaded configuration.
	svc := eks.NewFromConfig(cfg)

	// Call the DescribeCluster API to retrieve information about the specified cluster.
	resp, err := svc.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})
	if err != nil {
		// Return an error if an error occurs while calling the DescribeCluster API.
		return false, err
	}

	// Check if the cluster has an OIDC provider configured.
	if resp.Cluster.Identity.Oidc != nil {
		// Print a message indicating that the cluster has an OIDC provider configured.
		fmt.Println("Cluster has OIDC provider configured")

		// Return true to indicate that an OIDC provider exists for the cluster.
		return true, nil
	} else {
		// Print a message indicating that the cluster does not have an OIDC provider configured and provide instructions for configuring one.
		fmt.Println("Cluster does not have OIDC provider configured. Run 'eksctl utils associate-iam-oidc-provider --region="+region+" --cluster=", clusterName, " --approve' to associate an OIDC provider with the cluster")

		// Return false to indicate that an OIDC provider does not exist for the cluster.
		return false, nil
	}

}

// GetEksClusterNameByApiServerUrl is a function that retrieves the name of an Amazon EKS cluster using its API server URL.
// The function takes in the API server URL as an argument.
// The function returns the name of the cluster and an error if one occurs.
func GetEksClusterNameByApiServerUrl(apiServerUrl string) (string, error) {
	// Load the default AWS SDK configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Create a new Amazon EKS client using the loaded configuration.
	svc := eks.NewFromConfig(cfg)

	// List all Amazon EKS clusters in the current AWS account and region.
	resp, err := svc.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		// Return an error if an error occurs while listing the clusters.
		return "", err
	}

	// Loop through each cluster and compare its API server URL with the specified URL.
	for _, clusterName := range resp.Clusters {
		// Retrieve information about the current cluster.
		clusterInfo, err := svc.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: aws.String(clusterName),
		})
		if err != nil {
			// Return an error if an error occurs while describing the cluster.
			return "", err
		}

		// Check if the API server URL of the current cluster matches the specified URL.
		if *clusterInfo.Cluster.Endpoint == apiServerUrl {
			// Return the name of the current cluster if its API server URL matches the specified URL.
			return clusterName, nil
		}
	}

	// Return an error if no Amazon EKS cluster is found with the specified API server URL.
	return "", fmt.Errorf("could not find an EKS cluster with the API server URL %s", apiServerUrl)

}

// GetCurrentEKSClusterName is a function that retrieves the name of the Amazon EKS cluster currently in use by the kubectl command-line tool.
// The function first retrieves the API server endpoint of the current Kubernetes context by calling the GetCurrentKubeContextAPIEndpoint function.
// The function then retrieves the name of the Amazon EKS cluster associated with the API server endpoint by calling the GetEksClusterNameByApiServerUrl function.
// The function returns the name of the Amazon EKS cluster and an error if one occurs.
func GetCurrentEKSClusterName() (string, error) {
	// Retrieve the API server endpoint of the current Kubernetes context.
	apiEndpoint, err := GetCurrentKubeContextAPIEndpoint()
	if err != nil {
		// Return an error if an error occurs while retrieving the API server endpoint.
		return "", err
	}

	// Retrieve the name of the Amazon EKS cluster associated with the API server endpoint.
	clusterName, err := GetEksClusterNameByApiServerUrl(apiEndpoint)
	if err != nil {
		// Return an error if an error occurs while retrieving the Amazon EKS cluster name.
		return "", err
	}

	// Return the name of the Amazon EKS cluster.
	return clusterName, nil

}
