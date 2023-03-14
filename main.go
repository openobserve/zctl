/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/zinclabs/zctl/pkg/utils"
)

func main() {
	// cmd.Execute()

	releaseName := "zo1"
	utils.SetupAWS(releaseName)
	// utils.TearDownAWS(cluster, releaseName)

}
