package utils

import (
	"fmt"
)

// Setup function sets up AWS and Helm resources needed for the application.
// It takes the releaseName and namespace as input and returns an error if one occurs.
func Setup(installIdentifer, releaseName string, namespace string, region string) (SetupData, error) {

	setupData := SetupData{}

	// check if setup already exists
	config, err := ReadConfigMap("zincobserve-setup", namespace)
	if err == nil {
		fmt.Println("Setup already exists")
		fmt.Println(config)
		return setupData, err
	}

	bucket, role, err := SetupAWS(installIdentifer, releaseName, region)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return setupData, err
	}

	err = SetupHelm(releaseName, namespace, bucket, role)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up Helm resources.
		fmt.Println("error: ", err)
		return setupData, err
	}

	setupData = SetupData{
		Identifier:  installIdentifer,
		ReleaseName: releaseName,
		BucketName:  bucket,
		IamRole:     role,
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
		fmt.Println("error: ", err)
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
