package utils

import "fmt"

// SetupAWS sets up the necessary AWS resources for a given release.
// It returns the name of the S3 bucket and the IAM role ARN that were created.
// If an error occurs, it returns an empty string for both values and the error itself.
func SetupAWS(installIdentifer, releaseName string) (string, string, error) {
	// First, get the name of the current EKS cluster.
	clusterName, err := GetCurrentEKSClusterName()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	fmt.Println("EKS cluster name:", clusterName)

	// Next, get the default region for the AWS account.
	region, err := GetDefaultAwsRegion()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	// Set up the necessary AWS resources (S3 bucket and IAM role) for the release.
	bucketName, roleName, err := SetupAWSBase(installIdentifer, clusterName, releaseName, region)
	if err != nil {
		return "", "", err
	}

	// Return the names of the created resources.
	return bucketName, roleName, nil
}
