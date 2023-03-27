/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zinclabs/zctl/pkg/utils"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls a ZincObserve installation in the EKS cluster",
	Long: `
Uninstalls a ZincObserve installation in the EKS cluster. The subtasks include:
1. Uninstall the ConfigMap
2. Uninstall the IAM role
3. Uninstall the S3 bucket
4. Uninstall the helm release

	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()

		namespace := cmd.Flags().Lookup("namespace").Value.String()
		if namespace == "" {
			namespace, _ = utils.GetCurrentNamespace()
			fmt.Println("current namespace: ", namespace)
		}

		region := cmd.Flags().Lookup("region").Value.String()
		if region == "" {
			region, _ = utils.GetDefaultAwsRegion()
		}
		utils.Teardown(name, namespace, region)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	namespace := ""
	uninstallCmd.Flags().StringVar(&namespace, "namespace", "", "namespace to install the helm chart")

	region := ""
	uninstallCmd.Flags().StringVar(&region, "region", "", "region to delete the installation from ")
}
