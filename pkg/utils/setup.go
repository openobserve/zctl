package utils

import (
	"fmt"
)

// Setup function sets up AWS and Helm resources needed for the application.
// It takes the releaseName and namespace as input and returns an error if one occurs.
func Setup(installIdentifer, releaseName string, namespace string) error {
	bucket, roleArn, err := SetupAWS(installIdentifer, releaseName)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return err
	}
	err = SetupHelm(releaseName, namespace, bucket, roleArn)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up Helm resources.
		fmt.Println("error: ", err)
		return err
	}

	return nil
}

func Teardown(releaseName, namespace string) error {
	// Read the configmap
	cm, err := ReadConfigMap("zincobserve-setup", namespace)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return err
	}

	fmt.Println(cm["bucket_name"])

	err = TearDownAWS(releaseName, cm["bucket_name"], cm["role_arn"])
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return err
	}

	TearDownHelm(releaseName, namespace)

	return nil

}
