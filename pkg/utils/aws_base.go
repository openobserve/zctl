package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// SetupAWSBase creates an S3 bucket, IAM role and inline policy for the role. It returns the ARN of the role.
func SetupAWSBase(releaseIdentifer, clusterName, releaseName, region string) (string, string, error) {
	fmt.Println("..............Starting AWS Setup............")
	exists, err := HasOIDCProvider(clusterName, region)
	if err != nil {
		fmt.Println("error: ", err)
		return "", "", err
	}

	clusterDetails := &types.Cluster{}
	if exists {
		// Get EKS cluster details
		clusterDetails, err = GetEKSClusterDetails(clusterName)
		if err != nil {
			return "", "", err
		}
	}

	// capture items that we need in next steps
	issuer := *clusterDetails.Identity.Oidc.Issuer
	issuerId := issuer[len(issuer)-32:]
	awsAccountId, err := GetAWSAccountID()
	if err != nil {
		return "", "", err
	}

	// create an s3 bucket
	bucketName := "zinc-observe-" + releaseIdentifer + "-" + clusterName + "-" + releaseName
	err = CreateS3Bucket(bucketName)
	if err != nil {
		return "", "", err
	}

	// create an IAM role
	roleName := "zinc-observe-" + releaseIdentifer + "-" + clusterName + "-" + releaseName
	roleArn, err := CreateIAMRole(awsAccountId, region, issuerId, roleName, "zo-s3", clusterName, releaseName, bucketName)
	if err != nil {
		return "", "", err
	}

	return bucketName, roleArn, nil
}

// TearDownAWS tears down the AWS resources associated with a given release.
// It deletes the S3 bucket and the IAM role and policy.
// If an error occurs, it panics with the error message.
func TearDownAWS(releaseName string) error {
	releaseIdentifier := GetReleaseIdentifierFromReleaseName(releaseName)
	fmt.Println("..............TearDownAWS............")

	// First, get the name of the current EKS cluster.
	clusterName, err := GetCurrentEKSClusterName()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Delete the S3 bucket associated with the release.
	bucketName := "zinc-observe-" + releaseIdentifier + "-" + clusterName + "-" + releaseName
	err = DeleteS3Bucket(bucketName)
	if err != nil {
		return err
	}
	fmt.Println("S3 bucket deleted: ", bucketName)

	// Delete the IAM role and policy associated with the release.
	roleName := "zinc-observe-5080-" + releaseIdentifier + "-" + clusterName + "-" + releaseName
	err = DeleteIAMRoleWithPolicies(roleName)
	if err != nil {
		return err
	}

	return nil
}

// GetDefaultAwsRegion retrieves the default region for the AWS account.
// It returns the region string, or an error if one occurs.
func GetDefaultAwsRegion() (string, error) {
	// Load the AWS configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Get the default region from the configuration.
	region := cfg.Region

	return region, nil
}
