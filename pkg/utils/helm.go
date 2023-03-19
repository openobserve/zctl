package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// SetupHelm sets up the necessary kubernetes resources using official Helm chart.
// It returns the name of the S3 bucket and the IAM role ARN that were created.
// If an error occurs, it returns an empty string for both values and the error itself.
// It requires the name of the release, the namespace to deploy to, the name of the S3 bucket, and the IAM role ARN.
// If namespace is an empty string, it will default to "default". If namespace does not exist, it will be created.
func SetupHelm(releaseName, namespace, bucket, roleArn string) error {
	fmt.Println("install called")

	// Retrieve the URL of the Kubernetes cluster currently in use.
	clusterURL, err := GetCurrentKubeContextAPIEndpoint()
	if err != nil {
		// Print an error message if an error occurs while retrieving the cluster URL.
		fmt.Println("error: ", err)
		return err
	}

	// Retrieve the context of the Kubernetes cluster using its URL.
	context, err := KubeContextForCluster(clusterURL)
	if err != nil {
		// Print an error message if an error occurs while retrieving the context.
		fmt.Println("error: ", err)
		return err
	}

	// Create a new Helm object with the required deployment parameters.
	h1 := Helm{
		AppVersion:    "v0.2.0",
		ChartName:     "zincobserve",
		ChartVersion:  "0.2.6",
		Namespace:     namespace,
		ReleaseName:   releaseName,
		RepositoryURL: "https://charts.zinc.dev",
	}

	// Download the Helm chart specified by the Helm object.
	chart, err := h1.DownloadChart()
	if err != nil {
		// Print an error message if an error occurs while downloading the chart.
		fmt.Println("error downloading: ", err)
		return err
	}

	chart.Values, err = setUpChartValues(chart.Values, bucket, roleArn)
	if err != nil {
		// Print an error message if an error occurs while setting up the chart values.
		fmt.Println("error setting up chart values: ", err)
		return err
	}

	// Install the Helm chart with the updated values on the specified Kubernetes cluster context.
	err = h1.Install(chart, context)
	if err != nil {
		// Print an error message if an error occurs while installing the Helm chart.
		fmt.Println("error installing: ", err)
		return err
	}

	return nil

}

func TearDownHelm(releaseName, namespace string) {
	fmt.Println("TearDownHelm called")

	// Retrieve the URL of the Kubernetes cluster currently in use.
	// clusterURL, err := GetCurrentKubeContextAPIEndpoint()
	// if err != nil {
	// 	// Print an error message if an error occurs while retrieving the cluster URL.
	// 	fmt.Println("error: ", err)
	// }

	// Retrieve the context of the Kubernetes cluster using its URL.
	// context, err := KubeContextForCluster(clusterURL)
	// if err != nil {
	// 	// Print an error message if an error occurs while retrieving the context.
	// 	fmt.Println("error: ", err)
	// }

	// Create a new Helm object with the required deployment parameters.
	h1 := Helm{
		Namespace:   namespace,
		ReleaseName: releaseName,
	}

	// Uninstall the Helm chart on the specified Kubernetes cluster context.
	err := h1.UnInstall(releaseName, namespace)
	if err != nil {
		// Print an error message if an error occurs while uninstalling the Helm chart.
		fmt.Println("error uninstalling: ", err)
	}

}

func setUpChartValues(baseValuesMap map[string]interface{}, bucket, roleArn string) (map[string]interface{}, error) {
	// Marshal the values of the Helm chart to JSON format.
	jsonData, err := json.Marshal(baseValuesMap)
	if err != nil {
		// Print an error message if an error occurs while marshaling the values to JSON.
		fmt.Println("Error:", err)
		return nil, err
	}

	// Declare a variable to store the unmarshaled values from the Helm chart.
	var data ZincObserveValues

	// Unmarshal the values from JSON format and store them in the declared variable.
	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		// Print an error message if an error occurs while unmarshaling the values from JSON.
		fmt.Println("error unmarshalling: ", err)
		return nil, err
	}

	// Print a value from the unmarshaled data for testing purposes.
	fmt.Println(data.Auth.ZO_ROOT_USER_EMAIL)

	// Update the Helm chart values with the AWS bucket name and role ARN.
	data.Config.ZOS3BUCKETNAME = bucket
	data.ServiceAccount.Annotations["eks.amazonaws.com/role-arn"] = roleArn
	data.Image.Repository = "public.ecr.aws/zinclabs/zincobserve-dev"
	data.Image.Tag = "v0.2.0-42e58f1"

	yamlData, err := yaml.Marshal(&data)
	if err != nil {
		return nil, err
	}

	fmt.Println("YAML data: ", string(yamlData))

	// Convert the updated Helm chart values to a map and set them to the chart object.
	finalMap, err := StructToMap2(data)
	if err != nil {
		return nil, err
	}
	// fmt.Println("final Values.yaml: ")
	// PrettyPrint(finalMap)
	return finalMap, nil
}
