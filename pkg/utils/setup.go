package utils

import (
	"errors"
	"fmt"
)

// Setup function sets up AWS and Helm resources needed for the application.
// It takes the releaseName and namespace as input and returns an error if one occurs.
func Setup(inputData SetupData) (SetupData, error) {

	setupData := SetupData{}

	// check if setup already exists
	config, err := ReadConfigMap("zincobserve-setup", inputData.Namespace)
	if err == nil {
		fmt.Println("Setup already exists")
		fmt.Println(config)
		return setupData, err
	}

	if inputData.K8s == "eks" {

		bucket, role, clusterName, err := SetupAWS(inputData)
		if err != nil {
			// Print an error message and terminate the program if an error occurs while setting up AWS resources.
			fmt.Println("error: ", err)
			return setupData, err
		}

		inputData.BucketName = bucket
		inputData.IamRole = role
		inputData.ClusterName = clusterName

		accountId, err := GetAWSAccountID()
		if err != nil {
			fmt.Println("error: ", err)
			return setupData, err
		}
		roleArn := "arn:aws:iam::" + accountId + ":role/" + setupData.IamRole

		setupData = SetupData{
			Identifier:  inputData.Identifier,
			ReleaseName: inputData.ReleaseName,
			BucketName:  bucket,
			IamRole:     roleArn,
			K8s:         inputData.K8s,
			Region:      inputData.Region,
			ClusterName: clusterName,
		}

	} else if inputData.K8s == "gke" {
		// Setup GCP resources
		// 1. Get project ID
		// 2. Create a service account
		// 3. Create a bucket
		// 4. Create HMAC keys

		setup, err := SetupGCP(inputData)
		if err != nil {
			// Print an error message and terminate the program if an error occurs while setting up AWS resources.
			fmt.Println("error: ", err)
			return setupData, err
		}

		inputData.BucketName = setup.BucketName
		inputData.S3AccessKey = setup.S3AccessKey
		inputData.S3SecretKey = setup.S3SecretKey
		inputData.ServiceAccount = setup.ServiceAccount
		inputData.Region = "us-east-1"

		// setupData.BucketName = setup.BucketName
		// setupData.S3AccessKey = setup.S3AccessKey
		// setupData.S3SecretKey = setup.S3SecretKey
		// setupData.ServiceAccount = setup.ServiceAccount

	} else {
		return setupData, errors.New("k8s type not supported")
	}

	err = SetupHelm(inputData)
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

	err = TearDownAWS(cm, region)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return err
	}

	TearDownHelm(releaseName, namespace)

	DeleteConfigMap(cmName, namespace)

	return nil

}
