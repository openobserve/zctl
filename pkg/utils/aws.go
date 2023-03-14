package utils

import "fmt"

func SetupAWS(releaseName string) (string, string, error) {
	clusterName, err := GetCurrentEKSClusterName()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	fmt.Println("EKS cluster name:", clusterName)

	region, err := GetDefaultAwsRegion()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	bucketName, roleArn, err := SetupAWSBase(clusterName, releaseName, region)
	if err != nil {
		return "", "", err
	}

	return bucketName, roleArn, nil
}
