/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/zinclabs/zctl/pkg/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	// cmd.Execute()

	releaseName := "zo1"
	setupAWS(releaseName)
	// utils.TearDownAWS(cluster, releaseName)

}

func setupAWS(releaseName string) (string, error) {
	clusterName, err := utils.GetCurrentEKSClusterName()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("EKS cluster name:", clusterName)

	region, err := utils.GetDefaultAwsRegion()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	roleArn, err := utils.SetupAWS(clusterName, releaseName, region)
	if err != nil {
		return "", err
	}

	return roleArn, nil
}

// Create AWS IAM Role and Policy

func Install() {
	fmt.Println("install called")

	// clusterURL := utils.ClusterURLForCurrentContext()
	// context, err := utils.KubeContextForCluster(clusterURL)
	// if err != nil {
	// 	fmt.Println("error: ", err)
	// }

	h1 := utils.Helm{
		AppVersion:    "0.2.3",
		ChartName:     "zincobserve",
		ChartVersion:  "0.2.3",
		Namespace:     "t2",
		ReleaseName:   "zo2",
		RepositoryURL: "https://charts.zinc.dev",
	}

	chart, err := h1.DownloadChart()
	if err != nil {
		fmt.Println("error downloading: ", err)
	}

	// fmt.Println(chart.Values)

	// unmarshal chart.Values in map[string]interface{}

	jsonData, err := json.Marshal(chart.Values)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var data utils.ZincObserveValues

	err = yaml.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
	}

	fmt.Println(data.Auth.ZO_ROOT_USER_EMAIL)

	// err = h1.Install(chart, context)
	// if err != nil {
	// 	fmt.Println("error installing: ", err)
	// }
}
