package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func SetupAWS(clusterName, releaseName, region string) error {
	fmt.Println("..............Starting AWS Setup............")
	exists, err := HasOIDCProvider(clusterName, region)
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}

	clusterDetails := &types.Cluster{}
	if exists {
		// Get EKS cluster details
		clusterDetails, err = GetEKSClusterDetails(clusterName)
		if err != nil {
			return err
		}
		// print(clusterDetails)
	}

	// capture items that we need in next steps
	issuer := *clusterDetails.Identity.Oidc.Issuer
	issuerId := issuer[len(issuer)-32:]
	awsAccountId, err := GetAWSAccountID()
	if err != nil {
		return err
	}

	// create an s3 bucket
	bucketName := "zinc-observe-" + clusterName + "-" + releaseName
	err = CreateS3Bucket(bucketName)
	if err != nil {
		return err
	}

	// create an IAM role
	roleName := "zinc-observe-" + clusterName + "-" + releaseName
	err = CreateIAMRole(awsAccountId, region, issuerId, roleName, "zo-s3", clusterName, releaseName, bucketName)
	if err != nil {
		return err
	}

	return nil
}

func TearDownAWS(clusterName, releaseName string) {
	fmt.Println("..............TearDownAWS............")

	// Delete s3 bucket
	bucketName := "zinc-observe-" + clusterName + "-" + releaseName
	err := DeleteS3Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	// delete the IAM role and policy
	roleName := "zinc-observe-" + clusterName + "-" + releaseName
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
