package utils

import (
	"fmt"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetS3PolicyDocument(bucketName string) string {
	policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"s3:PutObject",
						"s3:GetObject",
						"s3:ListBucket",
						"s3:DeleteObject"
					],
					"Resource": [
						"arn:aws:s3:::%s",
						"arn:aws:s3:::%s/*"
					]
				}
			]
		}`, bucketName)

	return policy
}

// GetAWSAccountID retrieves the AWS account number for the current user.
// It returns the account number string, or an error if one occurs.
func GetAWSAccountID() (string, error) {

	// Load the AWS configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Create a new STS client to interact with AWS Security Token Service (STS).
	svc := sts.NewFromConfig(cfg)

	// Call the GetCallerIdentity API to retrieve the account number.
	resp, err := svc.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	// Extract the account number from the response.
	accountID := aws.ToString(resp.Account)

	return accountID, nil
}

// CreateIAMRole creates an IAM role with the EKS trusted entity and attaches an S3 bucket policy to it.
// It returns the ARN of the created role, or an error if one occurs.
func CreateIAMRole(accountId, region, issuerId, roleName, policyName, clusterName, releaseName, bucketName string) (string, error) {
	fmt.Println("Creating IAM role...")

	// Load the AWS configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	// Create a new IAM client.
	svc := iam.NewFromConfig(cfg)

	// Define the trusted entity for the role.
	trustedEntity := fmt.Sprintf(`{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Federated": "arn:aws:iam::%s:oidc-provider/oidc.eks.%s.amazonaws.com/id/%s"
			},
			"Action": "sts:AssumeRoleWithWebIdentity"
		}
	]
}`, accountId, region, issuerId)

	// Create the input for creating the role.
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(trustedEntity),
	}

	// Create the role.
	roleResp, err := svc.CreateRole(context.Background(), input)
	if err != nil {
		fmt.Println("Error in CreateRole: ", err.Error())
		return "", err
	}

	fmt.Println("Created IAM role: ", *roleResp.Role.Arn)

	fmt.Println("Creating inline policy for IAM role............")

	// Create a policy document for the S3 bucket policy.
	policyDocument := GetS3PolicyDocument(bucketName)

	// Attach the policy to the role.
	_, err = svc.PutRolePolicy(context.Background(), &iam.PutRolePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
		RoleName:       roleResp.Role.RoleName,
	})
	if err != nil {
		fmt.Println("Error in PutRolePolicy: ", err.Error())
		return "", err
	}

	fmt.Printf("Policy created for IAM role %s\n", roleName)
	return *roleResp.Role.Arn, nil
}

// DeleteIAMRoleWithPolicies deletes an IAM role and all of its associated policies.
// It returns an error if one occurs.
func DeleteIAMRoleWithPolicies(roleName string) error {
	fmt.Println("DeleteIAMRoleWithPolicies............")

	// Load the AWS configuration.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	// Create a new IAM client.
	svc := iam.NewFromConfig(cfg)

	// Delete any inline policies attached to the role.
	err = deleteInlineRolePolicies(roleName)
	if err != nil {
		return err
	}

	// Delete the role.
	_, err = svc.DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Deleted IAM role %s\n", roleName)
	return nil
}

// deleteInlineRolePolicies deletes all inline policies attached to an IAM role.
// It returns an error if one occurs.
func deleteInlineRolePolicies(roleName string) error {
	// Load the AWS SDK config.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	// Create a new IAM client.
	svc := iam.NewFromConfig(cfg)

	// List the inline policies attached to the role.
	resp, err := svc.ListRolePolicies(context.TODO(), &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	})
	if err != nil {
		return err
	}

	// Delete each inline policy attached to the role.
	for _, policyName := range resp.PolicyNames {
		_, err := svc.DeleteRolePolicy(context.TODO(), &iam.DeleteRolePolicyInput{
			RoleName:   &roleName,
			PolicyName: aws.String(policyName),
		})
		if err != nil {
			return err
		}
		fmt.Printf("Deleted inline policy %s from role %s\n", policyName, roleName)
	}

	return nil
}
