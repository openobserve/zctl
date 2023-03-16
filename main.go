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

	releaseName := "zo1"
	namespace := "default"
	// utils.Setup(releaseName, namespace)
	err := utils.SetupHelm(releaseName, namespace, "zinc-observe-5080-dev2-zo1", "arn:aws:iam::058694856476:role/zinc-observe-5080-dev2-zo1")
	if err != nil {
		// Print an error message and terminate the program if an error occurs while setting up Helm resources.
		fmt.Println("error: ", err)
		return
	}

	// utils.Teardown(releaseName, namespace)

}
