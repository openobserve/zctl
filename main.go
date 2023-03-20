/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/zinclabs/zctl/pkg/utils"
)

func main() {
	// cmd.Execute()

	// releaseName := "zo1"
	// namespace := "zo1"
	namespace, _ := utils.GetCurrentNamespace()
	// releaseIdentifer := utils.GenerateReleaseIdentifier()
	// utils.Setup(releaseIdentifer, releaseName, namespace)
	// err := utils.SetupHelm(releaseName, namespace, "zinc-observe-5080-dev2-zo1", "arn:aws:iam::058694856476:role/zinc-observe-5080-dev2-zo1")
	// if err != nil {
	// 	// Print an error message and terminate the program if an error occurs while setting up Helm resources.
	// 	fmt.Println("error: ", err)
	// 	return
	// }

	// utils.Teardown(releaseName, namespace)
	// utils.TearDownAWS(releaseName)

	name := "zincobserve-setup"

	// data := map[string]string{
	// 	"identifier":   "40075",
	// 	"release_name": releaseName,
	// 	"bucket_name":  "zinc-observe-40075-dev2-zo1",
	// 	"role_arn":     "arn:aws:iam::058694856476:role/zinc-observe-40075-dev2-zo1",
	// }

	// err := utils.CreateConfigMap(data, name, namespace)
	// if err != nil {
	// 	panic(err.Error())
	// }

	cm, err := utils.ReadConfigMap(name, namespace)
	if err != nil {
		panic(err.Error())
	}

	for key, value := range cm {
		fmt.Println(key, value)
	}
}
