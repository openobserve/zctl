package utils

import (
	"errors"
	"fmt"
)

// Setup function sets up AWS and Helm resources needed for the application.
// It takes the releaseName and namespace as input and returns an error if one occurs.
func Setup(setupData SetupData) (SetupData, error) {

	// check if setup already exists
	config, err := ReadConfigMap("zincobserve-setup", setupData.Namespace)
	if err == nil {
		fmt.Println("Setup already exists")
		fmt.Println(config)
		return setupData, err
	}

	if setupData.K8s == "eks" { ///////////////// Setup in EKS
		bucket, role, clusterName, err := SetupAWS(setupData)
		if err != nil {
			// Print an error message and terminate the program if an error occurs while setting up AWS resources.
			fmt.Println("error: ", err)
			return setupData, err
		}

		setupData.BucketName = bucket
		setupData.IamRole = role
		setupData.ClusterName = clusterName

		accountId, err := GetAWSAccountID()
		if err != nil {
			fmt.Println("error: ", err)
			return setupData, err
		}
		roleArn := "arn:aws:iam::" + accountId + ":role/" + setupData.IamRole

		setupData.BucketName = bucket
		setupData.IamRole = roleArn
		setupData.ClusterName = clusterName

	} else if setupData.K8s == "gke" { /////////////// Setup in GKE
		// Setup GCP resources
		// 1. Get project ID
		// 2. Create a service account
		// 3. Create a bucket
		// 4. Create HMAC keys

		gcpData, err := SetupGCP(setupData)
		if err != nil {
			// Print an error message and terminate the program if an error occurs while setting up AWS resources.
			fmt.Println("error: ", err)
			return setupData, err
		}

		setupData.BucketName = gcpData.BucketName
		setupData.S3AccessKey = gcpData.S3AccessKey
		setupData.S3SecretKey = gcpData.S3SecretKey
		setupData.ServiceAccount = gcpData.ServiceAccount
		setupData.Region = "us-east-1" // Dummy region required by aws sdk

	} else {
		return setupData, errors.New("k8s type not supported")
	}

	err = SetupHelm(setupData)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up Helm resources.
		fmt.Println("error: ", err)
		return setupData, err
	}

	return setupData, nil
}

func Teardown(releaseName, namespace, region string) error {

	// Get details from configmap
	cmName := "zincobserve-setup"
	// Read the configmap
	cm, err := ReadConfigMap(cmName, namespace)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error reading configmap for release: "+releaseName+" in namespace: "+namespace+" : ", err)
		return err
	}

	fmt.Println(cm)

	if cm.K8s == "eks" {
		err = TearDownAWS(cm, region)
		if err != nil {
			// Print an error message and terminate the program if an error occurs while setting up AWS resources.
			fmt.Println("error: ", err)
			return err
		}
	} else if cm.K8s == "gke" {
		// DeleteGCSBucket(cm) // We do not want to delete the data in the bucket
		DeleteGCPServiceAccount(cm)

	}

	TearDownHelm(releaseName, namespace)

	DeleteConfigMap(cmName, namespace)

	return nil

}
