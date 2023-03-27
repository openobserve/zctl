/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zinclabs/zctl/pkg/utils"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a ZincObserve installation in the EKS cluster",
	Long: `
Deletes a ZincObserve installation in the EKS cluster. The subtasks include:
1. Delete the ConfigMap
2. Delete the IAM role
3. Delete the S3 bucket
4. Delete the helm release

	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flags().Lookup("name").Value.String()

		namespace := cmd.Flags().Lookup("namespace").Value.String()
		if namespace == "" {
			namespace, _ = utils.GetCurrentNamespace()
		}
		utils.Teardown(name, namespace)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	namespace := ""
	deleteCmd.Flags().StringVar(&namespace, "namespace", "default", "namespace to install the helm chart")
	deleteCmd.MarkFlagRequired("namespace")
}
