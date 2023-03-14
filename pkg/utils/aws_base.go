package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// SetupAWSBase creates an S3 bucket, IAM role and inline policy for the role. It returns the ARN of the role.
func SetupAWSBase(clusterName, releaseName, region string) (string, string, error) {
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
	bucketName := "zinc-observe-5080-" + clusterName + "-" + releaseName
	err = CreateS3Bucket(bucketName)
	if err != nil {
		return "", "", err
	}

	// create an IAM role
	roleName := "zinc-observe-5080-" + clusterName + "-" + releaseName
	roleArn, err := CreateIAMRole(awsAccountId, region, issuerId, roleName, "zo-s3", clusterName, releaseName, bucketName)
	if err != nil {
		return "", "", err
	}

	return bucketName, roleArn, nil
}

func TearDownAWS(clusterName, releaseName string) {
	fmt.Println("..............TearDownAWS............")

	// Delete s3 bucket
	bucketName := "zinc-observe-5080-" + clusterName + "-" + releaseName
	err := DeleteS3Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	// delete the IAM role and policy
	roleName := "zinc-observe-5080-" + clusterName + "-" + releaseName
	err = DeleteIAMRoleWithPolicies(roleName)
	if err != nil {
		panic(err)
	}

}

func GetDefaultAwsRegion() (string, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Get the default region
	region := cfg.Region

	return region, nil
}
