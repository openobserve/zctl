package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

func InstallUsingHelm(releaseName string) {
	fmt.Println("install called")

	// Retrieve the URL of the Kubernetes cluster currently in use.
	clusterURL, err := GetCurrentKubeContextAPIEndpoint()
	if err != nil {
		// Print an error message if an error occurs while retrieving the cluster URL.
		fmt.Println("error: ", err)
		return
	}

	// Retrieve the context of the Kubernetes cluster using its URL.
	context, err := KubeContextForCluster(clusterURL)
	if err != nil {
		// Print an error message if an error occurs while retrieving the context.
		fmt.Println("error: ", err)
	}

	// Set up the required AWS resources for the application using a predefined setup function.
	bucket, roleArn, err := SetupAWS(releaseName)
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up AWS resources.
		fmt.Println("error: ", err)
		panic(err)
	}

	// Create a new Helm object with the required deployment parameters.
	h1 := Helm{
		AppVersion:    "0.2.3",
		ChartName:     "zincobserve",
		ChartVersion:  "0.2.3",
		Namespace:     "t2",
		ReleaseName:   releaseName,
		RepositoryURL: "https://charts.zinc.dev",
	}

	// Download the Helm chart specified by the Helm object.
	chart, err := h1.DownloadChart()
	if err != nil {
		// Print an error message if an error occurs while downloading the chart.
		fmt.Println("error downloading: ", err)
	}

	// Marshal the values of the Helm chart to JSON format.
	jsonData, err := json.Marshal(chart.Values)
	if err != nil {
		// Print an error message if an error occurs while marshaling the values to JSON.
		fmt.Println("Error:", err)
		return
	}

	// Declare a variable to store the unmarshaled values from the Helm chart.
	var data ZincObserveValues

	// Unmarshal the values from JSON format and store them in the declared variable.
	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		// Print an error message if an error occurs while unmarshaling the values from JSON.
		fmt.Println("error unmarshalling: ", err)
	}

	// Print a value from the unmarshaled data for testing purposes.
	fmt.Println(data.Auth.ZO_ROOT_USER_EMAIL)

	// Update the Helm chart values with the AWS bucket name and role ARN.
	data.Config.ZOS3BUCKETNAME = bucket
	data.ServiceAccount.Annotations["eks.amazonaws.com/role-arn"] = roleArn

	// Convert the updated Helm chart values to a map and set them to the chart object.
	chart.Values = StructToMap(data)

	// Install the Helm chart with the updated values on the specified Kubernetes cluster context.
	err = h1.Install(chart, context)
	if err != nil {
		// Print an error message if an error occurs while installing the Helm chart.
		fmt.Println("error installing: ", err)
	}

}
