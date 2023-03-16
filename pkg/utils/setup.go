package utils

import (
	"fmt"
)

// Setup function sets up AWS and Helm resources needed for the application.
// It takes the releaseName and namespace as input and returns an error if one occurs.
func Setup(releaseName string, namespace string) error {
	bucket, roleArn, err := SetupAWS(releaseName)
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
	err := TearDownAWS(releaseName)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		return err
	}

	TearDownHelm(releaseName, namespace)

	return nil

}
