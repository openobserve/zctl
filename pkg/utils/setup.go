package utils

import (
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

	bucket, role, clusterName, err := SetupAWS(inputData)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return setupData, err
	}

	inputData.BucketName = bucket
	inputData.IamRole = role
	inputData.ClusterName = clusterName

	err = SetupHelm(inputData)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up Helm resources.
		fmt.Println("error: ", err)
		return setupData, err
	}

	setupData = SetupData{
		Identifier:  inputData.Identifier,
		ReleaseName: inputData.ReleaseName,
		BucketName:  bucket,
		IamRole:     role,
		K8s:         inputData.K8s,
		Region:      inputData.Region,
		ClusterName: clusterName,
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
